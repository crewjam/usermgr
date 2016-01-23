package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type TestCatCommand struct {
	tempDir string
	Output  *bytes.Buffer
}

var _ = Suite(&TestCatCommand{})

func (s *TestCatCommand) SetUpTest(c *C) {
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

func (s *TestCatCommand) TearDownTest(c *C) {
	os.RemoveAll(s.tempDir)
}

func (s *TestCatCommand) TestCanPrintUser(c *C) {
	err := Main([]string{"usermgr", "--config=" + filepath.Join(s.tempDir, "usermgr.conf"),
		"cat", "alice"}, s.Output)
	c.Assert(err, IsNil)
	c.Assert(string(s.Output.Bytes()), Equals, ""+
		"{\n"+
		"  \"name\": \"alice\",\n"+
		"  \"real_name\": \"Alice Smith\",\n"+
		"  \"authorized_keys\": [\n"+
		"    \"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+ui4gptEr2ovoLD3vRhdRXXDLserFKhHcJrwBS79gO1J4KLzhgx0Pd/Mt7UyN3orxjKh06fd4N4P/5/c16BXK1Qe4DC/qClgkE5TyOyf8d04xXXVQlcn+LuRt4lAFgMxbfa2Sc0L0BJeu2VbW4DkIlYACwAdO6acWlOvJnMuYyomVgrcvle4yQWPU9L1Ql3E+RVIcdjR9aIN+QqgPNYZmvcuWzaKSbcnAwSsAIaoLxd8y14N6NvQdu4nvvZjBpkDTZI/IXIkwtZGkycSelNKnhPFWSL1qlgwqjH7U9/F3JxX4g0KjfzoCBjt9fKqn1fxneSZavFH1Q0LZNkfAUrov ross@rm\"\n"+
		"  ],\n"+
		"  \"backup_codes\": [\n"+
		"    {\n"+
		"      \"create_time\": \"2006-01-02T15:04:05Z\",\n"+
		"      \"salt\": \"QEJERkhKTE5QUlRWWFpcXmBiZGZoamxucHJ0dnh6fH4=\",\n"+
		"      \"hash\": \"P+yRmvoi7Gjw0vrdojCdZJlVqr6i5vUXzWuHxNXd7cOWvHWQG/ea8W/cs3YfVhcuECmozmhGPHYx/uku9/Impw==\"\n"+
		"    }\n"+
		"  ]\n"+
		"}")
}

func (s *TestCatCommand) TestCanPrintAll(c *C) {
	err := Main([]string{"usermgr", "--config=" + filepath.Join(s.tempDir, "usermgr.conf"),
		"cat"}, s.Output)
	c.Assert(err, IsNil)
	c.Assert(string(s.Output.Bytes()), Equals, ""+
		"{\n"+
		"  \"users\": [\n"+
		"    {\n"+
		"      \"name\": \"alice\",\n"+
		"      \"real_name\": \"Alice Smith\",\n"+
		"      \"authorized_keys\": [\n"+
		"        \"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+ui4gptEr2ovoLD3vRhdRXXDLserFKhHcJrwBS79gO1J4KLzhgx0Pd/Mt7UyN3orxjKh06fd4N4P/5/c16BXK1Qe4DC/qClgkE5TyOyf8d04xXXVQlcn+LuRt4lAFgMxbfa2Sc0L0BJeu2VbW4DkIlYACwAdO6acWlOvJnMuYyomVgrcvle4yQWPU9L1Ql3E+RVIcdjR9aIN+QqgPNYZmvcuWzaKSbcnAwSsAIaoLxd8y14N6NvQdu4nvvZjBpkDTZI/IXIkwtZGkycSelNKnhPFWSL1qlgwqjH7U9/F3JxX4g0KjfzoCBjt9fKqn1fxneSZavFH1Q0LZNkfAUrov ross@rm\"\n"+
		"      ],\n"+
		"      \"backup_codes\": [\n"+
		"        {\n"+
		"          \"create_time\": \"2006-01-02T15:04:05Z\",\n"+
		"          \"salt\": \"QEJERkhKTE5QUlRWWFpcXmBiZGZoamxucHJ0dnh6fH4=\",\n"+
		"          \"hash\": \"P+yRmvoi7Gjw0vrdojCdZJlVqr6i5vUXzWuHxNXd7cOWvHWQG/ea8W/cs3YfVhcuECmozmhGPHYx/uku9/Impw==\"\n"+
		"        }\n"+
		"      ]\n"+
		"    }\n"+
		"  ]\n"+
		"}")
}

func (s *TestCatCommand) TestInvalidUser(c *C) {
	err := Main([]string{"usermgr", "--config=" + filepath.Join(s.tempDir, "usermgr.conf"),
		"cat", "bob"}, s.Output)
	c.Assert(err, ErrorMatches, "bob: not found\n")
	c.Assert(string(s.Output.Bytes()), Equals, "")
}

func (s *TestCatCommand) TestUpdateFallsNoCache(c *C) {
	os.Remove(filepath.Join(s.tempDir, "users.pem"))
	err := Main([]string{"usermgr",
		"--config=" + filepath.Join(s.tempDir, "usermgr.conf"),
		"cat", "bob"}, s.Output)
	c.Assert(err, ErrorMatches, ".*: no such file or directory")
	c.Assert(string(s.Output.Bytes()), Equals, "")
}

func (s *TestCatCommand) TestBadConfig(c *C) {
	os.Remove(filepath.Join(s.tempDir, "usermgr.conf"))
	err := Main([]string{"usermgr",
		"--config=" + filepath.Join(s.tempDir, "usermgr.conf"),
		"cat", "bob"}, s.Output)
	c.Assert(err, ErrorMatches, ".*: no such file or directory")
	c.Assert(string(s.Output.Bytes()), Equals, "")
}
