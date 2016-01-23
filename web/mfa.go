package web

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/crewjam/usermgr"
	"github.com/pquerna/otp/totp"
	"golang.org/x/net/context"
)

func (s *Server) totpSecret(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	totpCode, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "not used",
		AccountName: "not used",
	})
	if err != nil {
		return err
	}

	now := TimeNow()
	device := usermgr.TOTPDevice{
		CreateTime: now,
	}
	device.SetSecret(s.AdminKey, totpCode.Secret())
	if err := device.GenerateCodes(now.Add(-10*time.Minute),
		now.Add(2*time.Hour), s.AdminKey); err != nil {
		return err
	}

	json.NewEncoder(w).Encode(struct {
		Secret string             `json:"secret"`
		Device usermgr.TOTPDevice `json:"device"`
	}{
		Secret: totpCode.Secret(),
		Device: device,
	})

	return nil
}

func (s *Server) makeBackupCode(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	buf := make([]byte, 10)
	if _, err := rand.Reader.Read(buf); err != nil {
		panic(err)
	}
	code := strings.ToLower(base32.StdEncoding.EncodeToString(buf))

	bc := usermgr.NewBackupCode(code)

	json.NewEncoder(w).Encode(struct {
		PlaintextCode string             `json:"plaintext_code"`
		BackupCode    usermgr.BackupCode `json:"backup_code"`
	}{
		PlaintextCode: code,
		BackupCode:    bc,
	})
	return nil
}
