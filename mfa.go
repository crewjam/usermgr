package usermgr

import (
	"fmt"
	"regexp"
	"time"

	"github.com/GeertJohan/yubigo"
)

var backupCode = regexp.MustCompile("^[A-Za-z0-9]{16}$")
var yubikeyCode = regexp.MustCompile("^[a-z]{44}$")
var totpCode = regexp.MustCompile("^[0-9]{6}$")

// ValidateCode checks that the specified code is a valid MFA code for the
// specified user. The supplied code can be one of the following:
//   - A TOTP code which is checked against validateURL
//   - A Yubikey code which is checked against validateURL
//   - A backup code which is compared to the hashed list of backup codes for the user
//
// Return nil if the code is valid, or an error otherwise.
func ValidateCode(user User, code string, yubicoClientID, yubicoSecretKey string) error {
	if yubikeyCode.MatchString(code) {
		yubiAuth, err := yubigo.NewYubiAuth(yubicoClientID, yubicoSecretKey)
		if err != nil {
			return err
		}

		yubikeyID := code[:12]
		for _, yubikeyDevice := range user.Yubikeys {
			if yubikeyDevice.DeviceID == yubikeyID {
				_, ok, err := yubiAuth.Verify(code)
				if err != nil {
					return fmt.Errorf("verify yubikey: %s", err)
				}
				if !ok {
					return fmt.Errorf("verify yubikey: code is not valid")
				}
				return nil
			}
		}
		return fmt.Errorf("verify yubikey: device is not registered")
	}

	if totpCode.MatchString(code) {
		for _, totpDevice := range user.TOTPDevices {
			err := totpDevice.VerifyCode(timeNow(), 30*time.Second, code)
			switch err {
			case nil:
				return nil // verified
			case ErrIncorrectCode:
				continue
			default:
				return fmt.Errorf("verify totp: %s", err)
			}
		}
		return fmt.Errorf("verify totp: %s", ErrIncorrectCode)
	}

	if backupCode.MatchString(code) {
		found := false
		for _, backupCode := range user.BackupCodes {
			if backupCode.Matches(code) {
				// TODO(ross): remove this code from the local cache so that it
				//   cannot be used again on this host.
				found = true
			}
		}
		if found {
			return nil
		}
	}

	return fmt.Errorf("invalid code")
}
