package web

import "golang.org/x/net/context"

type Storage interface {
	Get(ctx context.Context, etag string) (data []byte, newEtag string, err error)
	Put(ctx context.Context, data []byte) (etag string, err error)
}
