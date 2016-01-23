package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type TestSyncCommand struct {
	tempDir    string
	Output     *bytes.Buffer
	TestServer *httptest.Server
}

var _ = Suite(&TestSyncCommand{})

func (s *TestSyncCommand) SetUpTest(c *C) {
	s.tempDir, _ = ioutil.TempDir("", "unittest")

	s.TestServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	ioutil.WriteFile(filepath.Join(s.tempDir, "usermgr.conf"), []byte(""+
		fmt.Sprintf("CacheDir = %q\n", s.tempDir)+
		fmt.Sprintf("URL = %q\n", s.TestServer.URL)+
		"HostKey = \"m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8Pg\"\n"), 0644)

	s.Output = bytes.NewBuffer(nil)
}

func (s *TestSyncCommand) TearDownTest(c *C) {
	s.TestServer.Close()
	os.RemoveAll(s.tempDir)
}

func (s *TestSyncCommand) TestCanDoSyncOnce(c *C) {
	err := Main([]string{"usermgr",
		"--config=" + filepath.Join(s.tempDir, "usermgr.conf"),
		"sync"}, s.Output)

	c.Assert(err, IsNil)
	c.Assert(string(s.Output.Bytes()), Equals, "")
}

func (s *TestSyncCommand) TestUpdateFallsBackToCache(c *C) {
	s.TestServer.Close()
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
		"-----END USERMGR DATA-----\n"), 0644)

	err := Main([]string{"usermgr",
		"--config=" + filepath.Join(s.tempDir, "usermgr.conf"),
		"sync"}, s.Output)

	c.Assert(err, IsNil)
	c.Assert(string(s.Output.Bytes()), Equals, "")
}

func (s *TestSyncCommand) TestUpdateFallsNoCache(c *C) {
	s.TestServer.Close()

	err := Main([]string{"usermgr",
		"--config=" + filepath.Join(s.tempDir, "usermgr.conf"),
		"sync"}, s.Output)

	c.Assert(err, ErrorMatches, ".*: no such file or directory")
	c.Assert(string(s.Output.Bytes()), Equals, "")
}
func (s *TestSyncCommand) TestBadConfig(c *C) {
	os.Remove(filepath.Join(s.tempDir, "usermgr.conf"))

	err := Main([]string{"usermgr",
		"--config=" + filepath.Join(s.tempDir, "usermgr.conf"),
		"sync"}, s.Output)
	c.Assert(err, ErrorMatches, ".*: no such file or directory")
	c.Assert(string(s.Output.Bytes()), Equals, "")
}
