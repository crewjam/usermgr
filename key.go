package usermgr

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"

	"golang.org/x/crypto/nacl/box"
)

var randReader = rand.Reader

var ErrIncorrectKeyFormat = errors.New("incorrect key format")

type AdminKey struct {
	HostKey
	AdminPrivateKey [32]byte
	HostPublicKey   [32]byte
}

func (ak AdminKey) MarshalText() (text []byte, err error) {
	buf := make([]byte, 128)
	copy(buf[0:32], ak.AdminPublicKey[:])
	copy(buf[32:64], ak.HostPrivateKey[:])
	copy(buf[64:96], ak.AdminPrivateKey[:])
	copy(buf[96:128], ak.HostPublicKey[:])

	rv := make([]byte, base64.URLEncoding.EncodedLen(len(buf)))
	base64.URLEncoding.Encode(rv, buf)
	rv = rv[:len(rv)-1] // strip the trailing '='
	return rv, nil
}

func (ak *AdminKey) UnmarshalText(text []byte) error {
	text = append(text, '=')
	buf := make([]byte, base64.URLEncoding.DecodedLen(len(text)))
	n, err := base64.URLEncoding.Decode(buf, text)
	if err != nil || n != 128 {
		return ErrIncorrectKeyFormat
	}

	copy(ak.HostKey.AdminPublicKey[:], buf[0:32])
	buf = buf[32:]

	copy(ak.HostKey.HostPrivateKey[:], buf[0:32])
	buf = buf[32:]

	copy(ak.AdminPrivateKey[:], buf[0:32])
	buf = buf[32:]

	copy(ak.HostPublicKey[:], buf[0:32])
	return nil
}

func (ak AdminKey) MarshalJSON() ([]byte, error) {
	text, err := ak.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(text))
}

func (ak *AdminKey) UnmarshalJSON(b []byte) error {
	var text string
	if err := json.Unmarshal(b, &text); err != nil {
		return err
	}
	return ak.UnmarshalText([]byte(text))
}

func (ak AdminKey) String() string {
	text, _ := ak.MarshalText()
	return string(text)
}

type HostKey struct {
	AdminPublicKey [32]byte
	HostPrivateKey [32]byte
}

func (hk HostKey) MarshalText() (text []byte, err error) {
	buf := make([]byte, 64)
	copy(buf[0:32], hk.AdminPublicKey[:])
	copy(buf[32:64], hk.HostPrivateKey[:])
	rv := make([]byte, base64.URLEncoding.EncodedLen(64))
	base64.URLEncoding.Encode(rv, buf)
	rv = rv[:len(rv)-2] // strip the two trailing '='
	return rv, nil
}

func (hk *HostKey) UnmarshalText(text []byte) error {
	text = append(text, '=', '=')
	buf := make([]byte, base64.URLEncoding.DecodedLen(len(text)))
	n, err := base64.URLEncoding.Decode(buf, text)
	if err != nil || n != 64 {
		return ErrIncorrectKeyFormat
	}

	copy(hk.AdminPublicKey[:], buf[0:32])
	buf = buf[32:]

	copy(hk.HostPrivateKey[:], buf[0:32])
	return nil
}

func (hk HostKey) MarshalJSON() ([]byte, error) {
	text, err := hk.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(text))
}

func (hk *HostKey) UnmarshalJSON(b []byte) error {
	var text string
	if err := json.Unmarshal(b, &text); err != nil {
		return err
	}
	return hk.UnmarshalText([]byte(text))
}

func (hk HostKey) String() string {
	text, _ := hk.MarshalText()
	return string(text)
}

// GenerateKeyPair returnes a new key pair
func GenerateKeyPair() AdminKey {
	ak := AdminKey{}
	pub, priv, err := box.GenerateKey(randReader)
	if err != nil {
		panic(err)
	}
	ak.HostPublicKey, ak.HostPrivateKey = *pub, *priv

	pub, priv, err = box.GenerateKey(randReader)
	if err != nil {
		panic(err)
	}
	ak.AdminPublicKey, ak.AdminPrivateKey = *pub, *priv

	return ak
}
