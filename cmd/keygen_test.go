package cmd

import (
	"bytes"
	"os"
	"strings"

	. "gopkg.in/check.v1"
)

type TestKeygenCommand struct {
	tempDir string
	Output  *bytes.Buffer
}

var _ = Suite(&TestKeygenCommand{})

func (s *TestKeygenCommand) SetUpTest(c *C) {
	s.Output = bytes.NewBuffer(nil)
}

func (s *TestKeygenCommand) TearDownTest(c *C) {
	os.RemoveAll(s.tempDir)
}

func (s *TestKeygenCommand) TestCanPrintUser(c *C) {
	err := Main([]string{"usermgr", "keygen"}, s.Output)
	c.Assert(err, IsNil)

	lines := strings.Split(string(s.Output.Bytes()), "\n")
	c.Assert(lines[0], Matches, "admin key: [A-Za-z0-9_\\-]+")
	c.Assert(lines[1], Matches, "host key: [A-Za-z0-9_\\-]+")
}
