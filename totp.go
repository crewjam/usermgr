package usermgr

import (
	"crypto/hmac"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"time"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/scrypt"

	"github.com/pquerna/otp"
)

// TOTPDevice represents an enrolled TOTP device. Each device
// has a secret that is shared between us and the device which
// is used to generate time-based codes. The secret is stored
// here encrypted to the AdminKey, so it is only accessible to
// systems holding the admin key.
//
// To allow systems not holding the AdminKey to validate TOTP secrets
// a series of upcoming codes are stored in `Codes`. Since these
// codes are essentially passwords we can store them that way, using
// scrypt to generate a one-way hash which can be compared against
// what the user enters.
type TOTPDevice struct {
	Name       string    `json:"name,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`

	SecretNonce     []byte `json:"secret_nonce"`
	SecretEncrypted []byte `json:"secret_encrypted"`

	Codes []TOTPCode `json:"codes"`
}

func (d TOTPDevice) Secret(adminKey AdminKey) (string, error) {
	nonce := [24]byte{}
	copy(nonce[:], d.SecretNonce)

	plaintext, ok := box.Open(nil, d.SecretEncrypted, &nonce,
		&adminKey.AdminPublicKey, &adminKey.AdminPrivateKey)
	if !ok {
		return "", fmt.Errorf("cannot decrypt secret. Wrong key?")
	}
	return string(plaintext), nil
}

func (d *TOTPDevice) SetSecret(adminKey AdminKey, secret string) {
	nonce := [24]byte{}
	randReader.Read(nonce[:])
	d.SecretNonce = nonce[:]
	d.SecretEncrypted = box.Seal(nil, []byte(secret), &nonce,
		&adminKey.AdminPublicKey, &adminKey.AdminPrivateKey)
}

func (d *TOTPDevice) GenerateCodes(startTime, endTime time.Time, adminKey AdminKey) error {
	secret, err := d.Secret(adminKey)
	if err != nil {
		return err
	}

	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return fmt.Errorf("invalid secret: %s", err)
	}
	period := time.Second * 30

	// Compute scrypt hashes of TOTP codes in parallel.
	ch := make(chan TOTPCode)
	for t := startTime; t.Before(endTime); t = t.Add(period) {
		go func(t time.Time) {
			counter := uint64(math.Floor(float64(t.Unix()) / float64(period.Seconds())))
			code := ExpectedCode(counter, secretBytes)

			salt := make([]byte, saltBytes)
			_, err := io.ReadFull(randReader, salt)
			if err != nil {
				panic(err)
			}

			hash, err := scrypt.Key([]byte(code), salt, scryptN, scryptR, scryptP, hashBytes)
			if err != nil {
				panic(err)
			}

			ch <- TOTPCode{
				Time: t,
				Salt: salt,
				Hash: hash,
			}
		}(t)
	}
	for t := startTime; t.Before(endTime); t = t.Add(period) {
		d.Codes = append(d.Codes, <-ch)
	}
	return nil
}

var ErrIncorrectCode = errors.New("code is invalid")

func (d TOTPDevice) VerifyCode(now time.Time, skew time.Duration, userCode string) error {
	for _, code := range d.Codes {
		if now.Before(code.Time.Add(-1 * skew)) {
			continue
		}
		if now.After(code.Time.Add(skew)) {
			continue
		}
		actualHash, err := scrypt.Key([]byte(userCode), code.Salt, scryptN, scryptR, scryptP, hashBytes)
		if err != nil {
			return err
		}
		if subtle.ConstantTimeCompare(code.Hash, actualHash) == 1 {
			return nil // matched
		}
	}
	return ErrIncorrectCode
}

type TOTPCode struct {
	Time time.Time `json:"time"`
	Salt []byte    `json:"salt"`
	Hash []byte    `json:"hash"`
}

func ExpectedCode(counter uint64, secretBytes []byte) string {
	buf := make([]byte, 8)
	mac := hmac.New(otp.AlgorithmSHA1.Hash, secretBytes)
	binary.BigEndian.PutUint64(buf, counter)

	mac.Write(buf)
	sum := mac.Sum(nil)

	// "Dynamic truncation" in RFC 4226
	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	l := otp.DigitsSix.Length()
	mod := int32(value % int64(math.Pow10(l)))
	otpstr := otp.DigitsSix.Format(mod)
	return otpstr
}
