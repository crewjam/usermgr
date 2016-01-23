package usermgr

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/GeertJohan/yubigo"
	. "gopkg.in/check.v1"
)

var _ = Suite(&TestMFA{})

type TestMFA struct {
	AdminKey AdminKey
}

// RewriteTransport is an http.RoundTripper that rewrites requests
// using the provided URL's Scheme and Host, and its Path as a prefix.
// The Opaque field is untouched.
// If Transport is nil, http.DefaultTransport is used
type RewriteTransport struct {
	Transport http.RoundTripper
	URL       *url.URL
}

// RoundTrip implements the http.RoundTripper interface
func (t RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}

	if req.Header.Get("X-Amz-Date") != "" || strings.Contains(req.URL.String(), "amazonaws.com") {
		return rt.RoundTrip(req)
	}

	req.Header.Add("X-Original-Url", req.URL.String())
	// note that url.URL.ResolveReference doesn't work here
	// since t.u is an absolute url
	req.URL.Scheme = t.URL.Scheme
	req.URL.Host = t.URL.Host
	req.URL.Path = path.Join(t.URL.Path, req.URL.Path)
	return rt.RoundTrip(req)
}

func (s *TestMFA) SetUp(c *C) {
	s.AdminKey.UnmarshalText([]byte("m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8PkBCREZISkxOUFJUVlhaXF5gYmRmaGpsbnBydHZ4enx-ommQj5KJoeHRLhbHyA2RzNXBeJ_Xz4p1vJUsozZzhXw"))
}

func (s *TestMFA) TestValidate(c *C) {
	user := User{
		Name:     "alice",
		RealName: "Alice Smith",
		Yubikeys: []YubikeyDevice{YubikeyDevice{DeviceID: "vvnljckinlun"}},
		BackupCodes: []BackupCode{
			BackupCode{
				Salt: []byte{0x51, 0xe1, 0x15, 0x3f, 0x96, 0x59, 0x2e, 0x73, 0xab, 0x8d, 0xbf, 0x93, 0xc1, 0x2b, 0x20, 0xa7, 0x9c, 0x1c, 0x19, 0x46, 0x5e, 0x4c, 0xac, 0x5, 0x44, 0x4f, 0xc3, 0x26, 0x5f, 0xe9, 0xbb, 0x70},
				Hash: []byte{0x4c, 0x6, 0xd5, 0x49, 0x4, 0xfb, 0x52, 0xe, 0xa1, 0xe5, 0x4c, 0x74, 0xed, 0xf0, 0x2a, 0x62, 0x95, 0xd7, 0x7, 0xa2, 0x1f, 0xa6, 0x5f, 0x28, 0x7c, 0x3e, 0x2, 0x8c, 0xe6, 0xe7, 0x90, 0xce, 0x64, 0xdf, 0x39, 0xda, 0xd1, 0x2a, 0xf9, 0x7, 0xba, 0xe4, 0xd3, 0x2d, 0x26, 0xc, 0xad, 0xc8, 0xba, 0x9c, 0x2e, 0x8c, 0xd1, 0xa3, 0xc2, 0x36, 0x42, 0x99, 0xed, 0x37, 0xb3, 0x86, 0x8b, 0xb},
			},
		},
	}
	d := TOTPDevice{Name: "thing"}
	d.SetSecret(s.AdminKey, "WNZNSH53RSVUJWH2")
	d.GenerateCodes(time.Unix(1436983480, 0), time.Unix(1436983590, 0), s.AdminKey)
	user.TOTPDevices = []TOTPDevice{d}

	// yubikey
	func() {
		yubicoClientID := "12345"
		yubicoSecretKey := "" // an empty secret disables verification in the client
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Assert(r.Header.Get("X-Original-URL"), Matches,
				`^https://.*\.yubico.com/wsapi/2.0/verify\?id=12345&nonce=.*&otp=vvnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe&sl=secure$`)
			fmt.Fprintf(w, ""+
				"t=2014-07-15T18:26:43Z0888\r\n"+
				"otp=vvnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe\r\n"+
				"nonce=%s\r\n"+
				"sl=60\r\n"+
				"status=OK"+
				"\r\n"+
				"\r\n", r.FormValue("nonce"))
		}))
		defer testServer.Close()

		testServerURL, err := url.Parse(testServer.URL)
		c.Assert(err, IsNil)
		yubigo.HTTPClient = &http.Client{
			Transport: RewriteTransport{URL: testServerURL},
		}
		defer func() { yubigo.HTTPClient = nil }()

		err = ValidateCode(user, "vvnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe", yubicoClientID, yubicoSecretKey)
		c.Assert(err, IsNil)
	}()

	// yubikey bad secret
	func() {
		yubicoClientID := "12345"
		yubicoSecretKey := "wrongsekrit"
		err := ValidateCode(user, "vvnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe", yubicoClientID, yubicoSecretKey)
		c.Assert(err, ErrorMatches, "Given key seems to be invalid. Could not base64_decode. Error: illegal base64 data at input byte 8\n")
	}()

	// yubikey bad code
	func() {
		yubicoClientID := "12345"
		yubicoSecretKey := "" // an empty secret disables verification in the client
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Assert(r.Header.Get("X-Original-URL"), Matches,
				`^https://.*\.yubico.com/wsapi/2.0/verify\?id=12345&nonce=.*&otp=vvnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe&sl=secure$`)
			fmt.Fprintf(w, ""+
				"t=2014-07-15T18:26:43Z0888\r\n"+
				"otp=vvnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe\r\n"+
				"nonce=%s\r\n"+
				"sl=60\r\n"+
				"status=BAD_OTP"+
				"\r\n"+
				"\r\n", r.FormValue("nonce"))
		}))
		defer testServer.Close()

		testServerURL, err := url.Parse(testServer.URL)
		c.Assert(err, IsNil)
		yubigo.HTTPClient = &http.Client{
			Transport: RewriteTransport{URL: testServerURL},
		}
		defer func() { yubigo.HTTPClient = nil }()

		err = ValidateCode(user, "vvnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe", yubicoClientID, yubicoSecretKey)
		c.Assert(err, ErrorMatches, "verify yubikey: code is not valid")
	}()

	// yubikey fail
	func() {
		yubicoClientID := "12345"
		yubicoSecretKey := "" // an empty secret disables verification in the client
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Assert(r.Header.Get("X-Original-URL"), Matches,
				`^https://.*\.yubico.com/wsapi/2.0/verify\?id=12345&nonce=.*&otp=vvnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe&sl=secure$`)
			fmt.Fprintf(w, ""+
				"t=2014-07-15T18:26:43Z0888\r\n"+
				"otp=vvnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe\r\n"+
				"nonce=%s\r\n"+
				"sl=60\r\n"+
				"status=REPLAYED_OTP"+
				"\r\n"+
				"\r\n", r.FormValue("nonce"))
		}))
		defer testServer.Close()

		testServerURL, err := url.Parse(testServer.URL)
		c.Assert(err, IsNil)
		yubigo.HTTPClient = &http.Client{
			Transport: RewriteTransport{URL: testServerURL},
		}
		defer func() { yubigo.HTTPClient = nil }()

		err = ValidateCode(user, "vvnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe", yubicoClientID, yubicoSecretKey)
		c.Assert(err, ErrorMatches, "verify yubikey: The OTP is valid, but has been used before. If you receive this error, you might be the victim of a man-in-the-middle attack.")
	}()

	// yubikey unregistered
	func() {
		yubicoClientID := "12345"
		yubicoSecretKey := "" // an empty secret disables verification in the client
		err := ValidateCode(user, "xxnljckinlunnjknretikjnlgerkbkcvbhtcvgukubbe", yubicoClientID, yubicoSecretKey)
		c.Assert(err, ErrorMatches, "verify yubikey: device is not registered")
	}()

	// otp success
	func() {
		timeNow = func() time.Time { return time.Unix(1436983488, 0) }
		err := ValidateCode(user, "066243", "", "")
		c.Assert(err, IsNil)
	}()

	// expired
	func() {
		timeNow = func() time.Time { return time.Unix(1000000000, 0) }
		err := ValidateCode(user, "066243", "", "")
		c.Assert(err, ErrorMatches, "verify totp: code is invalid")
	}()

	// wrong code
	func() {
		timeNow = func() time.Time { return time.Unix(1436983488, 0) }
		err := ValidateCode(user, "123456", "", "")
		c.Assert(err, ErrorMatches, "verify totp: code is invalid")
	}()

	// wrong form
	func() {
		timeNow = func() time.Time { return time.Unix(1436983488, 0) }
		err := ValidateCode(user, "wrongformforcode", "", "")
		c.Assert(err, ErrorMatches, "invalid code")
	}()

	// backup code
	func() {
		err := ValidateCode(user, "Wyja8OeSygqwT9v8", "", "")
		c.Assert(err, IsNil)

		err = ValidateCode(user, "xxxxxxxxxxxxxxxx", "", "")
		c.Assert(err, ErrorMatches, "invalid code")
	}()
}
