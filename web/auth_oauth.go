package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/crewjam/httperr"
	"github.com/crewjam/usermgr"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const cookieMaxAge = time.Hour // TODO(ross): must be configurable
const cookieName = "token"

type OauthAuth struct {
	Config             oauth2.Config
	UserInfoURL        string
	EmailSuffix        string
	TokenSigningKey    []byte
	ValidateRemoteUser func(usermgr.User) (isAdmin bool, err error)
}

func (a OauthAuth) RequireUser(ctx context.Context, w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) { return []byte(a.TokenSigningKey), nil })
		if err == nil && token.Valid {
			return token.Claims["sub"].(string), nil
		}
		if err != nil {
			log.Printf("Cannot parse token: %s", err)
		} else if !token.Valid {
			log.Printf("token is not valid")
		}
	} else {
		log.Printf("cookie %s does not exist", cookieName)
	}

	state := jwt.New(jwt.SigningMethodHS256)
	state.Claims["url"] = r.URL.Path
	state.Claims["exp"] = TimeNow().Add(cookieMaxAge).Unix()
	stateString, err := state.SignedString(a.TokenSigningKey)
	if err != nil {
		return "", fmt.Errorf("cannot generate state JWT: %s", err)
	}
	http.Redirect(w, r, a.Config.AuthCodeURL(stateString), http.StatusFound)
	return "", httperr.Error{StatusCode: http.StatusFound}
}

func (a OauthAuth) HandleCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Verify the state
	state, err := jwt.Parse(r.FormValue("state"), func(token *jwt.Token) (interface{}, error) { return []byte(a.TokenSigningKey), nil })
	if err != nil {
		return err
	}
	if !state.Valid {
		return fmt.Errorf("state not valid")
	}

	// Exchange the code for a token
	oauthToken, err := a.Config.Exchange(ctx, r.FormValue("code"))
	if err != nil {
		return err
	}

	// Fetch the user info
	httpClient := a.Config.Client(ctx, oauthToken)
	userInfoResponse, err := httpClient.Get(a.UserInfoURL)
	if err != nil {
		return err
	}
	if userInfoResponse.StatusCode != http.StatusOK {
		return httperr.Error{
			StatusCode:   http.StatusForbidden,
			PrivateError: fmt.Errorf("user info: %s", userInfoResponse.Status),
		}
	}
	userInfo := OpenIDUserInfo{}
	err = json.NewDecoder(userInfoResponse.Body).Decode(&userInfo)
	if err != nil {
		return err
	}
	remoteUser := strings.TrimSuffix(userInfo.Email, a.EmailSuffix)

	// generate a token for the user
	jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))
	jwtToken.Claims["sub"] = remoteUser
	jwtToken.Claims["exp"] = TimeNow().Add(cookieMaxAge).Unix()
	tokenString, err := jwtToken.SignedString(a.TokenSigningKey)
	if err != nil {
		return err
	}

	// set the token and redirect
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Value:  tokenString,
		Path:   "/",
		Secure: true,
	})
	http.Redirect(w, r, state.Claims["url"].(string), http.StatusFound)
	return nil
}

// OpenIDUserInfo is the OpenID connect user information response
// http://openid.net/specs/openid-connect-core-1_0.html#UserInfoResponse
type OpenIDUserInfo struct {
	// Subject - Identifier for the End-User at the Issuer.
	Sub string `json:"sub"`

	// End-User's full name in displayable form including all name parts,
	// possibly including titles and suffixes, ordered according to the
	// End-User's locale and preferences.
	Name string `json:"name"`

	// End-User's preferred e-mail address. Its value MUST conform to the
	// RFC 5322 [RFC5322] addr-spec syntax. The RP MUST NOT rely upon this
	// value being unique, as discussed in Section 5.7.
	Email string `json:"email"`
}
