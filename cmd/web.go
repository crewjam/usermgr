package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/codegangsta/cli"
	"github.com/crewjam/usermgr/web"
	"github.com/zenazn/goji/web/middleware"
)

var webCommand = cli.Command{
	Name:   "web",
	Usage:  "Run the web interface",
	Action: WithError(WebCommand),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "bind, b",
			Value:  ":8000",
			Usage:  "Address to bind on.",
			EnvVar: "UM_BIND",
		},
		cli.StringFlag{
			Name:   "admin-key",
			Value:  "",
			Usage:  "The admin key",
			EnvVar: "UM_ADMIN_KEY",
		},
		cli.StringFlag{
			Name:   "store",
			Value:  "",
			Usage:  "The URL of the data storage service",
			EnvVar: "UM_STORE",
		},
		cli.StringFlag{
			Name:   "auth",
			Value:  "",
			Usage:  "The URL of the auth scheme. Example: oauth2://google?client_id=xxx&client_secret=yyy&email_suffix=@example.com",
			EnvVar: "UM_AUTH",
		},
		cli.StringFlag{
			Name:   "url",
			Value:  "",
			Usage:  "The root URL of the web service",
			EnvVar: "UM_URL",
		},
		cli.StringFlag{
			Name:   "token-key",
			Value:  "",
			Usage:  "the key used to sign auth tokens. should be random and secret",
			EnvVar: "UM_TOKEN_KEY",
		},
	},
}

func WebCommand(ctx *cli.Context) error {
	config := web.Config{}

	if err := config.AdminKey.UnmarshalText([]byte(ctx.String("admin-key"))); err != nil {
		return fmt.Errorf("cannot parse key: %s", err)
	}

	storeURL, err := url.Parse(ctx.String("store"))
	if err != nil {
		return fmt.Errorf("cannot parse store URL: %s", err)
	}
	switch storeURL.Scheme {
	case "file":
		config.Storage = web.FileStorage{Path: storeURL.Path}
	default:
		return fmt.Errorf("unknown scheme in store URL: %s", storeURL.String())
	}

	authURL, err := url.Parse(ctx.String("auth"))
	if err != nil {
		return fmt.Errorf("cannot parse auth URL: %s", err)
	}
	switch authURL.Scheme {
	case "oauth":
		scopes, ok := authURL.Query()["scope"]
		if !ok {
			scopes = []string{"openid", "profile", "email"}
		}

		endpoint := oauth2.Endpoint{
			AuthURL:  authURL.Query().Get("auth_url"),
			TokenURL: authURL.Query().Get("token_url"),
		}
		userInfoURL := authURL.Query().Get("user_info_url")

		if authURL.Host == "google" {
			endpoint = google.Endpoint
			userInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
		}

		config.Auth = web.OauthAuth{
			Config: oauth2.Config{
				ClientID:     authURL.Query().Get("client_id"),
				ClientSecret: authURL.Query().Get("client_secret"),
				Scopes:       scopes,
				RedirectURL:  fmt.Sprintf("%s/oauth2callback", ctx.String("url")),
				Endpoint:     endpoint,
			},
			UserInfoURL:     userInfoURL,
			EmailSuffix:     authURL.Query().Get("email_suffix"),
			TokenSigningKey: []byte(ctx.String("token-key")),
		}
	default:
		return fmt.Errorf("unknown scheme in auth URL: %s", authURL.String())
	}

	/*
		config := web.Config{

			Auth: web.OauthAuth{
				Config: oauth2.Config{
					ClientID:     "560985811258-9mu27s0raathl31lkr0j232hssikqhm2.apps.googleusercontent.com",
					ClientSecret: "ioJbOvdRFEiRA6B7jd5GT1Y5",
					Scopes:       []string{"openid", "profile", "email"},
					RedirectURL:  "https://96c42679.ngrok.io/oauth2callback",
					Endpoint:     google.Endpoint,
				},
				UserInfoURL:     "https://www.googleapis.com/oauth2/v3/userinfo",
				EmailSuffix:     "@octolabs.io",
				TokenSigningKey: []byte("XXX"),
			},
		}
	*/

	server := web.New(config)
	server.Mux.Use(middleware.RequestID)
	server.Mux.Use(middleware.Logger)

	log.Printf("listening on %s", ctx.String("bind"))
	return http.ListenAndServe(ctx.String("bind"), server.Mux)
}
