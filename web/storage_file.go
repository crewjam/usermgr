package web

import (
	"crypto/sha1"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
)

type FileStorage struct {
	Path string
}

func (fs FileStorage) Get(ctx context.Context, etag string) ([]byte, string, error) {
	existingEtagBuf, _ := ioutil.ReadFile(filepath.Join(fs.Path, "users.pem.etag"))
	if etag != "" && etag == string(existingEtagBuf) {
		return nil, etag, nil
	}

	buf, err := ioutil.ReadFile(filepath.Join(fs.Path, "users.pem"))
	if err != nil {
		return nil, "", err
	}

	return buf, string(existingEtagBuf), nil
}

func (fs FileStorage) Put(ctx context.Context, data []byte) (string, error) {
	etag := sha1.Sum(data)
	err := ioutil.WriteFile(filepath.Join(fs.Path, "users.pem.etag"), etag[:], 0644)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filepath.Join(fs.Path, "users.pem"), data, 0644)
	if err != nil {
		os.Remove(filepath.Join(fs.Path, "users.pem.etag"))
		return "", err
	}
	return string(etag[:]), err
}
