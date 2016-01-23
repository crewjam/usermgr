package appengine

import (
	"net/http"

	"github.com/crewjam/usermgr"
	"github.com/crewjam/usermgr/web"
	"google.golang.org/appengine"
)

func AdminKey(s string) usermgr.AdminKey {
	ak := usermgr.AdminKey{}
	if err := ak.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return ak
}

func init() {
	if config.Storage == nil {
		config.Storage = &Storage{}
	}
	web.NewContextFunc = appengine.NewContext

	server := web.New(config)
	http.Handle("/", server.Mux)
	http.Handle("/*", server.Mux)
}
