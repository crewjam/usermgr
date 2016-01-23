package usermgr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	. "gopkg.in/check.v1"
)

var _ = Suite(&TestLocal{})

type TestLocal struct {
	tempDir string
	HostKey HostKey
}

func (s *TestLocal) SetUpTest(c *C) {
	timeNow = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		return t
	}

	var err error
	s.tempDir, err = ioutil.TempDir("", "unittest")
	c.Assert(err, IsNil)
	s.HostKey.UnmarshalText([]byte("m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8Pg"))
}

func (s *TestLocal) TearDownTest(c *C) {
	os.RemoveAll(s.tempDir)
}

func (s *TestLocal) TestCanRead(c *C) {
	shouldFail := false
	testServerMissCount := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldFail {
			w.WriteHeader(http.StatusTeapot)
			return
		}
		if r.Header.Get("If-None-Match") == "this-is-the-etag" {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Add("ETag", "this-is-the-etag")
		fmt.Fprintf(w, ""+
			"-----BEGIN USERMGR DATA-----\n"+
			"AAIEBggKDA4QEhQWGBocHiAiJCYoKiwuZlypEp0ipVqsHQ8vZyddzC6DUFoU2hsH\n"+
			"iyN9qUK92GsFn03mkzm5JBety6Px9xCj7zdTbiIY+lPe/6AJJJc+hpDAzyfME5Zr\n"+
			"iu10rGNMPYegpMaSEraPJGeLBMXgANxdA+63mDxRQPKs9yp0e1Xd+SgrtBzd4ZV2\n"+
			"eBYJ5oe3nx10j2CAlPUdeBySOK9lpH8tU7T4U9leF2RHuL2aVsQdGPk7BMWve0a1\n"+
			"O+BLTKde8eas2SV4wCBu8inAjNftHDAopnhB+UwMBySI9l+1NG6CWL+oiSZDH3gz\n"+
			"TXfnVU8K71I99NOMYTyhsJUCxX2gjsquTxmZuPJ8R/Kjm2ckJXuMeg9nC+0P/Ilg\n"+
			"vAsMcOFCfpZDhEpYZ2p8OGbd/Vx8R16TLt+HBst3JPR97ron1B9JhjcHC05h8hRL\n"+
			"7L3YgcYdqYOjgf/qtUqm+gH3w1wfjwwrtAhvdxbqIwCgD5biriKmFbQpoqGpHpva\n"+
			"CnoATuwJQyqtp2QKP/it0ZiohidQ2in0hGVH5Wv9YlLGpznBt+RP1qL0g5Retirh\n"+
			"She72DsfQTxOIM6AKpvUA15vQO8Js57VCK1+0/CZ3nuo4RPyOCcv4qYVVbc+dCUT\n"+
			"POyA5kBQh9QmEdxJ7Wwa/2Kwupv/DxdZ3crjANPCUw/yxQs8v8YQAy2kRbDviud6\n"+
			"mH5NZDOrb8AowSsMaSv8c6eJYVkoShfzoC4JNxNYST3CvV1MK8MsoyEVhIw/N2F/\n"+
			"8eSKBmjr+aQfAX1QeY/DpoqMpivNeevd0dlByfVjV0PsiwixZezpqHt5FbncOfsD\n"+
			"x3xkAsRJhCPFOkKjbY2wAPZ+kRqi8V0IyHWPmn7YD3pllnlO+iH4d+sap7rnkeNa\n"+
			"1Y1UtjureHSTrpYWV06uoFGdWrgBx2yyzEs4HAfe0IoWRZBmNBCqPQo=\n"+
			"-----END USERMGR DATA-----\n")
		testServerMissCount++
	}))
	defer testServer.Close()

	expectedUserData := &UsersData{Users: []User{
		User{
			Name:           "alice",
			RealName:       "Alice Smith",
			AuthorizedKeys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+ui4gptEr2ovoLD3vRhdRXXDLserFKhHcJrwBS79gO1J4KLzhgx0Pd/Mt7UyN3orxjKh06fd4N4P/5/c16BXK1Qe4DC/qClgkE5TyOyf8d04xXXVQlcn+LuRt4lAFgMxbfa2Sc0L0BJeu2VbW4DkIlYACwAdO6acWlOvJnMuYyomVgrcvle4yQWPU9L1Ql3E+RVIcdjR9aIN+QqgPNYZmvcuWzaKSbcnAwSsAIaoLxd8y14N6NvQdu4nvvZjBpkDTZI/IXIkwtZGkycSelNKnhPFWSL1qlgwqjH7U9/F3JxX4g0KjfzoCBjt9fKqn1fxneSZavFH1Q0LZNkfAUrov ross@rm"},
			BackupCodes: []BackupCode{BackupCode{
				// "hellohellohelloo"
				CreateTime: timeNow(),
				Salt:       []byte{0x40, 0x42, 0x44, 0x46, 0x48, 0x4a, 0x4c, 0x4e, 0x50, 0x52, 0x54, 0x56, 0x58, 0x5a, 0x5c, 0x5e, 0x60, 0x62, 0x64, 0x66, 0x68, 0x6a, 0x6c, 0x6e, 0x70, 0x72, 0x74, 0x76, 0x78, 0x7a, 0x7c, 0x7e},
				Hash:       []byte{0x3f, 0xec, 0x91, 0x9a, 0xfa, 0x22, 0xec, 0x68, 0xf0, 0xd2, 0xfa, 0xdd, 0xa2, 0x30, 0x9d, 0x64, 0x99, 0x55, 0xaa, 0xbe, 0xa2, 0xe6, 0xf5, 0x17, 0xcd, 0x6b, 0x87, 0xc4, 0xd5, 0xdd, 0xed, 0xc3, 0x96, 0xbc, 0x75, 0x90, 0x1b, 0xf7, 0x9a, 0xf1, 0x6f, 0xdc, 0xb3, 0x76, 0x1f, 0x56, 0x17, 0x2e, 0x10, 0x29, 0xa8, 0xce, 0x68, 0x46, 0x3c, 0x76, 0x31, 0xfe, 0xe9, 0x2e, 0xf7, 0xf2, 0x26, 0xa7},
			}},
		},
	}}

	ud, err := UpdateLocalCache(s.tempDir, testServer.URL, s.HostKey)
	c.Assert(err, IsNil)
	c.Assert(ud, DeepEquals, expectedUserData)
	c.Assert(testServerMissCount, Equals, 1)

	buf, err := ioutil.ReadFile(filepath.Join(s.tempDir, "users.pem"))
	c.Assert(err, IsNil)
	c.Assert(string(buf), DeepEquals, ""+
		"-----BEGIN USERMGR DATA-----\n"+
		"AAIEBggKDA4QEhQWGBocHiAiJCYoKiwuZlypEp0ipVqsHQ8vZyddzC6DUFoU2hsH\n"+
		"iyN9qUK92GsFn03mkzm5JBety6Px9xCj7zdTbiIY+lPe/6AJJJc+hpDAzyfME5Zr\n"+
		"iu10rGNMPYegpMaSEraPJGeLBMXgANxdA+63mDxRQPKs9yp0e1Xd+SgrtBzd4ZV2\n"+
		"eBYJ5oe3nx10j2CAlPUdeBySOK9lpH8tU7T4U9leF2RHuL2aVsQdGPk7BMWve0a1\n"+
		"O+BLTKde8eas2SV4wCBu8inAjNftHDAopnhB+UwMBySI9l+1NG6CWL+oiSZDH3gz\n"+
		"TXfnVU8K71I99NOMYTyhsJUCxX2gjsquTxmZuPJ8R/Kjm2ckJXuMeg9nC+0P/Ilg\n"+
		"vAsMcOFCfpZDhEpYZ2p8OGbd/Vx8R16TLt+HBst3JPR97ron1B9JhjcHC05h8hRL\n"+
		"7L3YgcYdqYOjgf/qtUqm+gH3w1wfjwwrtAhvdxbqIwCgD5biriKmFbQpoqGpHpva\n"+
		"CnoATuwJQyqtp2QKP/it0ZiohidQ2in0hGVH5Wv9YlLGpznBt+RP1qL0g5Retirh\n"+
		"She72DsfQTxOIM6AKpvUA15vQO8Js57VCK1+0/CZ3nuo4RPyOCcv4qYVVbc+dCUT\n"+
		"POyA5kBQh9QmEdxJ7Wwa/2Kwupv/DxdZ3crjANPCUw/yxQs8v8YQAy2kRbDviud6\n"+
		"mH5NZDOrb8AowSsMaSv8c6eJYVkoShfzoC4JNxNYST3CvV1MK8MsoyEVhIw/N2F/\n"+
		"8eSKBmjr+aQfAX1QeY/DpoqMpivNeevd0dlByfVjV0PsiwixZezpqHt5FbncOfsD\n"+
		"x3xkAsRJhCPFOkKjbY2wAPZ+kRqi8V0IyHWPmn7YD3pllnlO+iH4d+sap7rnkeNa\n"+
		"1Y1UtjureHSTrpYWV06uoFGdWrgBx2yyzEs4HAfe0IoWRZBmNBCqPQo=\n"+
		"-----END USERMGR DATA-----\n")
	buf, err = ioutil.ReadFile(filepath.Join(s.tempDir, "users.pem.etag"))
	c.Assert(err, IsNil)
	c.Assert(string(buf), DeepEquals, "this-is-the-etag")

	// now try with an existing etag
	testServerMissCount = 0
	ud, err = UpdateLocalCache(s.tempDir, testServer.URL, s.HostKey)
	c.Assert(err, IsNil)
	c.Assert(ud, DeepEquals, expectedUserData)
	c.Assert(testServerMissCount, Equals, 0)

	// if the cached data is invalid, we don't use it
	testServerMissCount = 0
	ioutil.WriteFile(filepath.Join(s.tempDir, "users.pem"), []byte("wrong-content"), 0644)
	ud, err = UpdateLocalCache(s.tempDir, testServer.URL, s.HostKey)
	c.Assert(err, IsNil)
	c.Assert(ud, DeepEquals, expectedUserData)
	c.Assert(testServerMissCount, Equals, 1)

	// try again with an old etag
	testServerMissCount = 0
	ioutil.WriteFile(filepath.Join(s.tempDir, "users.pem.etag"), []byte("wrong-etag"), 0644)
	ud, err = UpdateLocalCache(s.tempDir, testServer.URL, s.HostKey)
	c.Assert(err, IsNil)
	c.Assert(ud, DeepEquals, expectedUserData)
	c.Assert(testServerMissCount, Equals, 1)

	// server fails
	testServerMissCount = 0
	shouldFail = true
	ud, err = UpdateLocalCache(s.tempDir, testServer.URL, s.HostKey)
	c.Assert(err, ErrorMatches, "418 I'm a teapot")
	c.Assert(ud, IsNil)
	c.Assert(testServerMissCount, Equals, 0)
	shouldFail = false
}

func (s *TestLocal) TestCannotReadEtagFile(c *C) {
	_, err := UpdateLocalCache("/dev/null", "", s.HostKey)
	c.Assert(err, ErrorMatches, "open /dev/null/users.pem.etag: not a directory")
}

func (s *TestLocal) TestInvalidURLs(c *C) {
	_, err := UpdateLocalCache(s.tempDir, "h_t_t_p://foo%3", s.HostKey)
	c.Assert(err, ErrorMatches, "parse h_t_t_p://foo%3: invalid URL escape \"%3\"")

	_, err = UpdateLocalCache(s.tempDir, "bar://foo", s.HostKey)
	c.Assert(err, ErrorMatches, "Get bar://foo: unsupported protocol scheme \"bar\"")
}

func (s *TestLocal) TestInvalidData(c *C) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("ETag", "this-is-the-etag")
		fmt.Fprintf(w, "this is not valid data")
	}))
	defer testServer.Close()

	_, err := UpdateLocalCache(s.tempDir, testServer.URL, s.HostKey)
	c.Assert(err, ErrorMatches, "invalid encoding")

	_, err = ioutil.ReadFile(filepath.Join(s.tempDir, "users.pem.etag"))
	c.Assert(os.IsNotExist(err), Equals, true)
	_, err = ioutil.ReadFile(filepath.Join(s.tempDir, "users.pem"))
	c.Assert(os.IsNotExist(err), Equals, true)
}

func (s *TestLocal) TestCannotWriteToCacheDir(c *C) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("ETag", "this-is-the-etag")
		fmt.Fprintf(w, ""+
			"-----BEGIN USERMGR DATA-----\n"+
			"AAIEBggKDA4QEhQWGBocHiAiJCYoKiwuZlypEp0ipVqsHQ8vZyddzC6DUFoU2hsH\n"+
			"iyN9qUK92GsFn03mkzm5JBety6Px9xCj7zdTbiIY+lPe/6AJJJc+hpDAzyfME5Zr\n"+
			"iu10rGNMPYegpMaSEraPJGeLBMXgANxdA+63mDxRQPKs9yp0e1Xd+SgrtBzd4ZV2\n"+
			"eBYJ5oe3nx10j2CAlPUdeBySOK9lpH8tU7T4U9leF2RHuL2aVsQdGPk7BMWve0a1\n"+
			"O+BLTKde8eas2SV4wCBu8inAjNftHDAopnhB+UwMBySI9l+1NG6CWL+oiSZDH3gz\n"+
			"TXfnVU8K71I99NOMYTyhsJUCxX2gjsquTxmZuPJ8R/Kjm2ckJXuMeg9nC+0P/Ilg\n"+
			"vAsMcOFCfpZDhEpYZ2p8OGbd/Vx8R16TLt+HBst3JPR97ron1B9JhjcHC05h8hRL\n"+
			"7L3YgcYdqYOjgf/qtUqm+gH3w1wfjwwrtAhvdxbqIwCgD5biriKmFbQpoqGpHpva\n"+
			"CnoATuwJQyqtp2QKP/it0ZiohidQ2in0hGVH5Wv9YlLGpznBt+RP1qL0g5Retirh\n"+
			"She72DsfQTxOIM6AKpvUA15vQO8Js57VCK1+0/CZ3nuo4RPyOCcv4qYVVbc+dCUT\n"+
			"POyA5kBQh9QmEdxJ7Wwa/2Kwupv/DxdZ3crjANPCUw/yxQs8v8YQAy2kRbDviud6\n"+
			"mH5NZDOrb8AowSsMaSv8c6eJYVkoShfzoC4JNxNYST3CvV1MK8MsoyEVhIw/N2F/\n"+
			"8eSKBmjr+aQfAX1QeY/DpoqMpivNeevd0dlByfVjV0PsiwixZezpqHt5FbncOfsD\n"+
			"x3xkAsRJhCPFOkKjbY2wAPZ+kRqi8V0IyHWPmn7YD3pllnlO+iH4d+sap7rnkeNa\n"+
			"1Y1UtjureHSTrpYWV06uoFGdWrgBx2yyzEs4HAfe0IoWRZBmNBCqPQo=\n"+
			"-----END USERMGR DATA-----\n")
	}))
	defer testServer.Close()

	func() {
		os.Chmod(s.tempDir, 0500)
		defer os.Chmod(s.tempDir, 0700)
		_, err := UpdateLocalCache(s.tempDir, testServer.URL, s.HostKey)
		c.Assert(err, ErrorMatches, ".* permission denied")
	}()

	func() {
		ioutil.WriteFile(filepath.Join(s.tempDir, "users.pem.etag"), []byte("xxx"), 0400)
		defer os.Chmod(filepath.Join(s.tempDir, "users.pem.etag"), 0600)

		_, err := UpdateLocalCache(s.tempDir, testServer.URL, s.HostKey)
		c.Assert(err, ErrorMatches, ".* permission denied")

		_, err = ioutil.ReadFile(filepath.Join(s.tempDir, "users.pem~"))
		c.Assert(os.IsNotExist(err), Equals, true)
	}()

	func() {
		ioutil.WriteFile(filepath.Join(s.tempDir, "users.pem.etag"), []byte("xxx"), 0600)
		ioutil.WriteFile(filepath.Join(s.tempDir, "users.pem~"), []byte("xxx"), 0600)
		ioutil.WriteFile(filepath.Join(s.tempDir, "users.pem"), []byte("xxx"), 0400)
		os.Chmod(s.tempDir, 0500)
		defer os.Chmod(s.tempDir, 0700)
		defer os.Chmod(filepath.Join(s.tempDir, "users.pem"), 0600)

		_, err := UpdateLocalCache(s.tempDir, testServer.URL, s.HostKey)
		c.Assert(err, ErrorMatches, "rename .*: permission denied")

		// Note: users.pem.etag and user.pem~ are not deleted because of
		// the directory permissions, which means we cannot test that the
		// cleanup works in the event renaming users.pem~ to users.pem fails.
	}()
}

func (s *TestLocal) TestGetLocalCache(c *C) {
	ioutil.WriteFile(filepath.Join(s.tempDir, "users.pem"), []byte(""+
		"-----BEGIN USERMGR DATA-----\n"+
		"AAIEBggKDA4QEhQWGBocHiAiJCYoKiwuZlypEp0ipVqsHQ8vZyddzC6DUFoU2hsH\n"+
		"iyN9qUK92GsFn03mkzm5JBety6Px9xCj7zdTbiIY+lPe/6AJJJc+hpDAzyfME5Zr\n"+
		"iu10rGNMPYegpMaSEraPJGeLBMXgANxdA+63mDxRQPKs9yp0e1Xd+SgrtBzd4ZV2\n"+
		"eBYJ5oe3nx10j2CAlPUdeBySOK9lpH8tU7T4U9leF2RHuL2aVsQdGPk7BMWve0a1\n"+
		"O+BLTKde8eas2SV4wCBu8inAjNftHDAopnhB+UwMBySI9l+1NG6CWL+oiSZDH3gz\n"+
		"TXfnVU8K71I99NOMYTyhsJUCxX2gjsquTxmZuPJ8R/Kjm2ckJXuMeg9nC+0P/Ilg\n"+
		"vAsMcOFCfpZDhEpYZ2p8OGbd/Vx8R16TLt+HBst3JPR97ron1B9JhjcHC05h8hRL\n"+
		"7L3YgcYdqYOjgf/qtUqm+gH3w1wfjwwrtAhvdxbqIwCgD5biriKmFbQpoqGpHpva\n"+
		"CnoATuwJQyqtp2QKP/it0ZiohidQ2in0hGVH5Wv9YlLGpznBt+RP1qL0g5Retirh\n"+
		"She72DsfQTxOIM6AKpvUA15vQO8Js57VCK1+0/CZ3nuo4RPyOCcv4qYVVbc+dCUT\n"+
		"POyA5kBQh9QmEdxJ7Wwa/2Kwupv/DxdZ3crjANPCUw/yxQs8v8YQAy2kRbDviud6\n"+
		"mH5NZDOrb8AowSsMaSv8c6eJYVkoShfzoC4JNxNYST3CvV1MK8MsoyEVhIw/N2F/\n"+
		"8eSKBmjr+aQfAX1QeY/DpoqMpivNeevd0dlByfVjV0PsiwixZezpqHt5FbncOfsD\n"+
		"x3xkAsRJhCPFOkKjbY2wAPZ+kRqi8V0IyHWPmn7YD3pllnlO+iH4d+sap7rnkeNa\n"+
		"1Y1UtjureHSTrpYWV06uoFGdWrgBx2yyzEs4HAfe0IoWRZBmNBCqPQo=\n"+
		"-----END USERMGR DATA-----\n"), 0400)
	ud, err := GetLocalCache(s.tempDir, s.HostKey)
	c.Assert(err, IsNil)
	c.Assert(ud.Users[0].Name, Equals, "alice")
}

func (s *TestLocal) TestGetLocalCacheFail(c *C) {
	ioutil.WriteFile(filepath.Join(s.tempDir, "users.pem"), []byte("invalid"), 0400)
	_, err := GetLocalCache(s.tempDir, s.HostKey)
	c.Assert(err, ErrorMatches, "invalid encoding")

	os.Chmod(filepath.Join(s.tempDir, "users.pem"), 0)
	defer os.Chmod(filepath.Join(s.tempDir, "users.pem"), 0644)
	_, err = GetLocalCache(s.tempDir, s.HostKey)
	c.Assert(err, ErrorMatches, ".*: permission denied")
}
