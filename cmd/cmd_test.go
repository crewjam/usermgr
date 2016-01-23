package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/crewjam/usermgr"

	. "gopkg.in/check.v1"
)

type TestMainCommand struct {
	tempDir string
	Output  *bytes.Buffer
}

var _ = Suite(&TestMainCommand{})

func (s *TestMainCommand) SetUpTest(c *C) {
	s.tempDir, _ = ioutil.TempDir("", "unittest")
	ioutil.WriteFile(filepath.Join(s.tempDir, "usermgr.conf"), []byte(""+
		fmt.Sprintf("CacheDir = %q\n", s.tempDir)+
		"HostKey = \"m_NiqMyWkkgOi1sT4uMCnp5kYuNanescRkRr3DP29FUAAgQGCAoMDhASFBYYGhweICIkJigqLC4wMjQ2ODo8Pg\"\n"), 0644)
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
	s.Output = bytes.NewBuffer(nil)
}

func (s *TestMainCommand) TearDownTest(c *C) {
	os.RemoveAll(s.tempDir)
}

func (s *TestMainCommand) TestCanPrintUser(c *C) {
	config, err := LoadConfig(filepath.Join(s.tempDir, "usermgr.conf"))
	c.Assert(err, IsNil)
	c.Assert(config, DeepEquals, &Config{
		URL: "",
		HostKey: usermgr.HostKey{
			AdminPublicKey: [32]uint8{0x9b, 0xf3, 0x62, 0xa8, 0xcc, 0x96, 0x92, 0x48, 0xe, 0x8b, 0x5b, 0x13, 0xe2, 0xe3, 0x2, 0x9e, 0x9e, 0x64, 0x62, 0xe3, 0x5a, 0x9d, 0xeb, 0x1c, 0x46, 0x44, 0x6b, 0xdc, 0x33, 0xf6, 0xf4, 0x55},
			HostPrivateKey: [32]uint8{0x0, 0x2, 0x4, 0x6, 0x8, 0xa, 0xc, 0xe, 0x10, 0x12, 0x14, 0x16, 0x18, 0x1a, 0x1c, 0x1e, 0x20, 0x22, 0x24, 0x26, 0x28, 0x2a, 0x2c, 0x2e, 0x30, 0x32, 0x34, 0x36, 0x38, 0x3a, 0x3c, 0x3e},
		},
		CacheDir:         s.tempDir,
		LoginGroups:      []string{"users"},
		SudoGroups:       []string{"wheel"},
		LoginMFARequried: false,
	})

	config, err = LoadConfig(filepath.Join(s.tempDir, "missing"))
	c.Assert(err, ErrorMatches, ".*: no such file or directory")
	c.Assert(config, IsNil)

	os.Chmod(filepath.Join(s.tempDir, "usermgr.conf"), 0)
	defer os.Chmod(filepath.Join(s.tempDir, "usermgr.conf"), 0644)
	config, err = LoadConfig(filepath.Join(s.tempDir, "usermgr.conf"))
	c.Assert(err, ErrorMatches, ".*: permission denied")
	c.Assert(config, IsNil)

	ioutil.WriteFile(filepath.Join(s.tempDir, "broken.conf"), []byte("CacheDir = \n"), 0644)
	config, err = LoadConfig(filepath.Join(s.tempDir, "broken.conf"))
	c.Assert(err, ErrorMatches, "Near line 1 .* Expected value but found .*")
	c.Assert(config, IsNil)
}
