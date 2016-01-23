package appengine

import (
	"crypto/sha1"
	"fmt"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type StoredData struct {
	Data []byte
	Etag string
}

type Storage struct {
}

func (Storage) Get(ctx context.Context, existingEtag string) ([]byte, string, error) {
	storedData := StoredData{}
	key := datastore.NewKey(ctx, "StoredData", "stored_data", 0, nil)
	err := datastore.Get(ctx, key, &storedData)
	if err == datastore.ErrNoSuchEntity {
		return nil, "", os.ErrNotExist
	}
	if err != nil {
		return nil, "", err
	}
	if storedData.Etag == existingEtag {
		return nil, existingEtag, nil
	}
	return storedData.Data, storedData.Etag, nil
}

func (Storage) Put(ctx context.Context, data []byte) (string, error) {
	key := datastore.NewKey(ctx, "StoredData", "stored_data", 0, nil)
	etag := fmt.Sprintf("%x", sha1.Sum(data))
	_, err := datastore.Put(ctx, key, &StoredData{Data: data, Etag: etag})
	if err != nil {
		return "", err
	}
	return etag, nil
}
