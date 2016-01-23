package usermgr

import (
	"crypto/subtle"
	"io"
	"time"

	"golang.org/x/crypto/scrypt"
)

const (
	scryptN   = 16384
	scryptR   = 8
	scryptP   = 1
	saltBytes = 32
	hashBytes = 64
)

var timeNow = time.Now

type BackupCode struct {
	CreateTime time.Time `json:"create_time,omitempty"`
	Salt       []byte    `json:"salt,omitempty"`
	Hash       []byte    `json:"hash,omitempty"`
}

func NewBackupCode(code string) BackupCode {
	salt := make([]byte, saltBytes)
	_, err := io.ReadFull(randReader, salt)
	if err != nil {
		panic(err)
	}

	hash, err := scrypt.Key([]byte(code), salt, scryptN, scryptR, scryptP, hashBytes)
	if err != nil {
		panic(err)
	}

	return BackupCode{
		CreateTime: timeNow(),
		Salt:       salt,
		Hash:       hash,
	}
}

func (bc BackupCode) Matches(userCode string) bool {
	if len(bc.Salt) != saltBytes || len(bc.Hash) != hashBytes {
		return false
	}

	actualHash, err := scrypt.Key([]byte(userCode), bc.Salt, scryptN, scryptR, scryptP, hashBytes)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(bc.Hash, actualHash) == 1
}
