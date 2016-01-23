// package web implemented the web interface for usermgr
package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/crewjam/httperr"
	"github.com/crewjam/usermgr"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
)

var TimeNow = time.Now

type Config struct {
	Storage     Storage
	Auth        Auth
	AdminKey    usermgr.AdminKey
	DownloadURL string
}

type RemoteUser struct {
	Name    string
	IsAdmin bool
}

type Server struct {
	Mux         *web.Mux
	Storage     Storage
	Auth        Auth
	ContextFunc func() context.Context
	AdminKey    usermgr.AdminKey
	DownloadURL string
}

func (s *Server) getSignedData(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	etag := r.Header.Get("If-None-Match")
	data, newEtag, err := s.Storage.Get(ctx, etag)
	if err != nil {
		return err
	}
	if data == nil {
		w.WriteHeader(http.StatusNotModified)
		return nil
	}
	w.Header().Set("ETag", newEtag)
	w.Write(data)
	return nil
}

func (s *Server) mutateUsersData(ctx context.Context, f func(usersData *usermgr.UsersData) error) error {
	usersData, _, err := s.loadData(ctx)
	if err != nil {
		return err
	}

	if err := f(usersData); err != nil {
		return err
	}

	if err := s.storeData(ctx, usersData); err != nil {
		return err
	}

	return nil
}

func (s *Server) cronHourly(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if err := s.mutateUsersData(ctx, func(usersData *usermgr.UsersData) error {
		startTime := TimeNow().Add(-10 * time.Minute)
		endTime := TimeNow().Add(2 * time.Hour)
		for _, user := range usersData.Users {
			for i, device := range user.TOTPDevices {
				if err := device.GenerateCodes(startTime, endTime, s.AdminKey); err != nil {
					return err
				}
				user.TOTPDevices[i] = device
			}
			usersData.Set(user)
		}
		return nil
	}); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) postGlobal(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	remoteUser, err := s.RequireUser(ctx, w, r)
	if err != nil {
		return err
	}
	if !remoteUser.IsAdmin {
		return httperr.Forbidden
	}

	if err := s.mutateUsersData(ctx, func(usersData *usermgr.UsersData) error {
		if k := r.FormValue("yubikey_client_id"); k != "" {
			usersData.YubikeyClientID = k
		}
		if k := r.FormValue("yubikey_client_secret"); k != "" {
			usersData.YubikeyClientSecret = k
		}
		return nil
	}); err != nil {
		return err
	}
	return s.getIndex(ctx, w, r)
}

func (s *Server) getIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	remoteUser, err := s.RequireUser(ctx, w, r)
	if err != nil {
		return err
	}

	buf, _ := index_html()
	w.Write(buf)
	fmt.Fprintf(w, "<script>hostKey=\"%s\";</script>", s.AdminKey.HostKey)
	if remoteUser.IsAdmin {
		w.Write([]byte("<script>isAdmin=true;</script>"))
	}
	return nil
}

func (s *Server) loadData(ctx context.Context) (*usermgr.UsersData, string, error) {
	usersDataBuf, etag, err := s.Storage.Get(ctx, "")
	if os.IsNotExist(err) {
		return &usermgr.UsersData{}, "", nil
	}
	if err != nil {
		return nil, "", err
	}

	usersData, err := usermgr.LoadUsersData(usersDataBuf, s.AdminKey.HostKey)
	if err != nil {
		return nil, "", err
	}
	return usersData, etag, nil
}

func (s *Server) storeData(ctx context.Context, usersData *usermgr.UsersData) error {
	signedUserData, err := usersData.SignedString(s.AdminKey)
	if err != nil {
		return err
	}

	if _, err := s.Storage.Put(ctx, signedUserData); err != nil {
		return err
	}
	return nil
}

type contextKeyType int

const (
	urlParamsKey contextKeyType = iota
)

// Param returns the named URL parameter or an empty string if it is not present
func Param(ctx context.Context, name string) string {
	rv, _ := ctx.Value(urlParamsKey).(map[string]string)[name]
	return rv
}

// NewContextFunc is invoked to renerate the new context for
// a request. The default implementation uses context.Background,
// but when running in appengine this must be replaced with
// appengine.NewContext(r).
var NewContextFunc = func(r *http.Request) context.Context {
	return context.Background()
}

type wrapRequest func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func (f wrapRequest) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ctx := NewContextFunc(r)
	ctx = context.WithValue(ctx, urlParamsKey, c.URLParams)
	err := f(ctx, w, r)
	if err != nil {
		log.Printf("ERROR: %s", err)
		httperr.Write(w, err)
	}
}

func New(config Config) *Server {
	s := Server{
		Mux:         web.New(),
		Storage:     config.Storage,
		Auth:        config.Auth,
		ContextFunc: func() context.Context { return context.Background() },
		AdminKey:    config.AdminKey,
		DownloadURL: config.DownloadURL,
	}
	if s.DownloadURL == "" {
		s.DownloadURL = "https://github.com/crewjam/usermgr/releases/download/XXX/usermgr"
	}

	s.Mux.Get("/setup", wrapRequest(s.getSetup))
	s.Mux.Get("/setup/:key", wrapRequest(s.getSetup))
	s.Mux.Get("/users.pem", wrapRequest(s.getSignedData))
	s.Mux.Post("/_totp", wrapRequest(s.totpSecret))
	s.Mux.Post("/_backup_code", wrapRequest(s.makeBackupCode))
	s.Mux.Post("/_cron/hourly", wrapRequest(s.cronHourly))

	s.Mux.Get("/", wrapRequest(s.getIndex))
	s.Mux.Post("/", wrapRequest(s.postGlobal))
	s.Mux.Get("/users/", wrapRequest(s.getUsersList))
	s.Mux.Get("/users/:user", wrapRequest(s.getUser))
	s.Mux.Put("/users/:user", wrapRequest(s.putUser))
	s.Mux.Delete("/users/:user", wrapRequest(s.deleteUser))
	if oauth, ok := s.Auth.(OauthAuth); ok {
		s.Mux.Get("/oauth2callback", wrapRequest(oauth.HandleCallback))
	}

	return &s
}
