package web

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"golang.org/x/net/context"

	. "gopkg.in/check.v1"

	"github.com/crewjam/httperr"
	"github.com/crewjam/usermgr"
)

var _ = Suite(&TestWeb{})

type TestWeb struct {
	FakeStorage *FakeStorage
	FakeAuth    *FakeAuth
	AdminKey    usermgr.AdminKey
	Server      *Server
}

type FakeStorage struct {
	Data []byte
	Etag string
	Err  error
}

func (fs *FakeStorage) Get(ctx context.Context, etag string) ([]byte, string, error) {
	if etag == fs.Etag {
		return nil, fs.Etag, fs.Err
	}
	return fs.Data, fs.Etag, fs.Err
}

func (fs *FakeStorage) Put(ctx context.Context, data []byte) (string, error) {
	fs.Data = data
	fs.Etag = fmt.Sprintf("%x", sha1.Sum(data))
	return fs.Etag, fs.Err
}

type FakeAuth struct {
	User string
	Err  error
}

func (fa *FakeAuth) RequireUser(ctx context.Context, w http.ResponseWriter, r *http.Request) (string, error) {
	return fa.User, fa.Err
}

func (suite *TestWeb) SetUpTest(c *C) {
	suite.FakeStorage = &FakeStorage{}
	suite.FakeStorage.Put(context.TODO(), []byte(""+
		"-----BEGIN USERMGR DATA-----\n"+
		"n5yAQIwaoyBk4bYUzb0fvtIAej9sTDsX8Zwvoatb6vm9o6zEyIx7YPEwNYPDu+/6\n"+
		"dLcFVxuQ9IbPhczlrIuWZf8d3hvrCWPNXptYEkt6P8h553gxXH7xIAVC6sidYn4u\n"+
		"8/pWezCFPa8W9gNz80sQV2pLDQiuOYk+HkvNgb5KKyLWtyoY9cNAu2Ms866cLXHI\n"+
		"KyXBdv1Lz44EaDEovIqn8aotkKNx18vLzReAiw8IEGCR63leu87JEdoSpt1gdMWx\n"+
		"XZAuJTH8xOvHb2b+auK5m4CUglw0RQMhKZ1N9+NfxFesB7SQqyTj46clFK3H5+Cu\n"+
		"ZQ723urcpXv3U0hfqLrPdagr5qtqTdP7jI/eUgyJrX4lUBuSStlraTJRrhM8xouf\n"+
		"4TPWHo6vLfKmjYm0SDmFr6nesYLq7EGmjX34CtcprKlzKEZJctmu77adrnS9lIly\n"+
		"sWlx8mIPy6VaC50T9/0IlA6EOv/PS2AyfurxaexrpZDNieUaAuda0estjsqhWJdJ\n"+
		"fAFRrOmjmtPIAonLkfzP0S4OPEMRfbfZ4zeXsPjJZmioh2dz30cgsKDMptxhPUGP\n"+
		"5PZAT5i4CVrT0a6hfP8SXROqEx7IgpGevtYmCiV2iCd5pahdtI4JnDRcRbfPvlzY\n"+
		"17KIrTQ9JMurWlQeglozJru91fkgF+VnHclw3Eq+RW+UbAArppHLHDRkkBvPctTM\n"+
		"As8DQPU+txmFp2kwqfoLP9rLxvi280T4KAPMtUeYbHBFizp7ymJwlmnDDcbzNjtP\n"+
		"bTuzCpXAAUYD5w5SJnGVTbq0gsWf7TTaIHEcLuc/36tL2Oar8TA8YDJeqhQGBQPc\n"+
		"LJEtE5VVbflw7JJHdzkrZQQIcKNO7jG8w7/s5rak6i2ANuu2jUHjFMokELcMWD/M\n"+
		"s3CTLY+rLTOaVvfaE/frs7p6gK5NWsdXaIvFhQS5TZh7ElgyYIu07zMeXen3C/pG\n"+
		"gZwoGXcArFSHyGHWYtD3+B1no+2+oPIj/P5MOp8Nk26yi/8WtTU+J2bYEloNMSPa\n"+
		"rV6hfH0P\n"+
		"-----END USERMGR DATA-----\n"))
	suite.AdminKey.UnmarshalText([]byte("m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXw"))

	suite.FakeAuth = &FakeAuth{
		User: "alice",
		Err:  nil,
	}

	suite.Server = New(Config{
		Storage:  suite.FakeStorage,
		Auth:     suite.FakeAuth,
		AdminKey: suite.AdminKey,
	})

}

func (suite *TestWeb) TearDownTest(c *C) {
}

func (suite *TestWeb) TestCanGetSignedData(c *C) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users.pem", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 200)
	c.Assert(string(w.Body.Bytes()), Equals, ""+
		"-----BEGIN USERMGR DATA-----\n"+
		"n5yAQIwaoyBk4bYUzb0fvtIAej9sTDsX8Zwvoatb6vm9o6zEyIx7YPEwNYPDu+/6\n"+
		"dLcFVxuQ9IbPhczlrIuWZf8d3hvrCWPNXptYEkt6P8h553gxXH7xIAVC6sidYn4u\n"+
		"8/pWezCFPa8W9gNz80sQV2pLDQiuOYk+HkvNgb5KKyLWtyoY9cNAu2Ms866cLXHI\n"+
		"KyXBdv1Lz44EaDEovIqn8aotkKNx18vLzReAiw8IEGCR63leu87JEdoSpt1gdMWx\n"+
		"XZAuJTH8xOvHb2b+auK5m4CUglw0RQMhKZ1N9+NfxFesB7SQqyTj46clFK3H5+Cu\n"+
		"ZQ723urcpXv3U0hfqLrPdagr5qtqTdP7jI/eUgyJrX4lUBuSStlraTJRrhM8xouf\n"+
		"4TPWHo6vLfKmjYm0SDmFr6nesYLq7EGmjX34CtcprKlzKEZJctmu77adrnS9lIly\n"+
		"sWlx8mIPy6VaC50T9/0IlA6EOv/PS2AyfurxaexrpZDNieUaAuda0estjsqhWJdJ\n"+
		"fAFRrOmjmtPIAonLkfzP0S4OPEMRfbfZ4zeXsPjJZmioh2dz30cgsKDMptxhPUGP\n"+
		"5PZAT5i4CVrT0a6hfP8SXROqEx7IgpGevtYmCiV2iCd5pahdtI4JnDRcRbfPvlzY\n"+
		"17KIrTQ9JMurWlQeglozJru91fkgF+VnHclw3Eq+RW+UbAArppHLHDRkkBvPctTM\n"+
		"As8DQPU+txmFp2kwqfoLP9rLxvi280T4KAPMtUeYbHBFizp7ymJwlmnDDcbzNjtP\n"+
		"bTuzCpXAAUYD5w5SJnGVTbq0gsWf7TTaIHEcLuc/36tL2Oar8TA8YDJeqhQGBQPc\n"+
		"LJEtE5VVbflw7JJHdzkrZQQIcKNO7jG8w7/s5rak6i2ANuu2jUHjFMokELcMWD/M\n"+
		"s3CTLY+rLTOaVvfaE/frs7p6gK5NWsdXaIvFhQS5TZh7ElgyYIu07zMeXen3C/pG\n"+
		"gZwoGXcArFSHyGHWYtD3+B1no+2+oPIj/P5MOp8Nk26yi/8WtTU+J2bYEloNMSPa\n"+
		"rV6hfH0P\n"+
		"-----END USERMGR DATA-----\n")
	c.Assert(w.Header().Get("Etag"), DeepEquals, "34cfb4411e1ed35d1183e544202c0608e3c91c0c")

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/users.pem", nil)
	r.Header.Add("If-None-Match", "34cfb4411e1ed35d1183e544202c0608e3c91c0c")
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusNotModified)
	c.Assert(string(w.Body.Bytes()), Equals, "")

	// storage fail
	suite.FakeStorage.Data = nil
	suite.FakeStorage.Etag = ""
	suite.FakeStorage.Err = fmt.Errorf("cannot frob the grob")
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/users.pem", nil)
	r.Header.Add("If-None-Match", "34cfb4411e1ed35d1183e544202c0608e3c91c0c")
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusInternalServerError)
	c.Assert(string(w.Body.Bytes()), Equals, "Internal Server Error\n")
}

func (suite *TestWeb) TestCanFetchIndex(c *C) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	suite.Server.Mux.ServeHTTP(w, r)

	c.Assert(w.Code, Equals, 200)
}

func (suite *TestWeb) TestCron(c *C) {
	suite.FakeAuth.Err = fmt.Errorf("not reached")

	// add a TOTP token to alice
	ud, _ := usermgr.LoadUsersData(suite.FakeStorage.Data, suite.AdminKey.HostKey)
	alice := ud.Users[0]
	dev := usermgr.TOTPDevice{}
	dev.SetSecret(suite.AdminKey, "WNZNSH53RSVUJWH2")
	alice.TOTPDevices = append(alice.TOTPDevices, dev)
	ud.Set(alice)
	suite.FakeStorage.Data, _ = ud.SignedString(suite.AdminKey)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/_cron/hourly", nil)
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 204)

	ud, _ = usermgr.LoadUsersData(suite.FakeStorage.Data, suite.AdminKey.HostKey)
	c.Assert(len(ud.Users[0].TOTPDevices[0].Codes), Equals, 261)
}

func (suite *TestWeb) TestGlobal(c *C) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/", strings.NewReader("yubikey_client_id=one&yubikey_client_secret=password"))
	r.Header.Set("Content-type", "application/x-www-form-urlencoded")
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, 200)

	ud, _ := usermgr.LoadUsersData(suite.FakeStorage.Data, suite.AdminKey.HostKey)
	c.Assert(ud.YubikeyClientID, Equals, "one")
	c.Assert(ud.YubikeyClientSecret, Equals, "password")

	suite.FakeAuth.User = "bob"
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/", strings.NewReader("yubikey_client_id=one&yubikey_client_secret=password"))
	r.Header.Set("Content-type", "application/x-www-form-urlencoded")
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusForbidden)

	suite.FakeAuth.Err = httperr.Unauthorized
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/", strings.NewReader("yubikey_client_id=one&yubikey_client_secret=password"))
	r.Header.Set("Content-type", "application/x-www-form-urlencoded")
	suite.Server.Mux.ServeHTTP(w, r)
	c.Assert(w.Code, Equals, http.StatusUnauthorized)
}
