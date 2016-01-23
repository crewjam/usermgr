package web

import (
	"net/http"

	"golang.org/x/net/context"
)

type NilAuth struct {
}

func (na NilAuth) RequireUser(ctx context.Context, w http.ResponseWriter, r *http.Request) (string, error) {
	return "anonymous", nil
}
