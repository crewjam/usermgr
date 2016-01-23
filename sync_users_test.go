package usermgr

import (
	"fmt"
	"io/ioutil"
	"os"

	. "gopkg.in/check.v1"
)

var _ = Suite(&TestUsers{})

type ExecCall struct {
	Name string
	Args []string
}

type TestUsers struct {
	tempDir     string
	LoginGroups []string
	DB          *UsersData

	ExecCalls []ExecCall
}

func (s *TestUsers) SetUpTest(c *C) {
	s.LoginGroups = []string{"global-user", "myapp-user"}
	s.DB = &UsersData{
		Users: []User{
			// alice exists and remains
			{
				Name:   "alice",
				Groups: []string{"myapp-user"},
			},
			// bob is deleted because he doesn't exist
			// charlie is deleted because he isn't in myapp-use
			{
				Name:   "charlie",
				Groups: []string{"myapp-admin"},
			},
			// badluckbrian cannot be created but we try
			{
				Name:   "badluckbrian",
				Groups: []string{"myapp-user"},
			},
			// dave doesn't exist and is created
			{
				Name:   "dave",
				Groups: []string{"myapp-user"},
			},
		},
	}

	var err error
	s.tempDir, err = ioutil.TempDir("", "unittest")
	c.Assert(err, IsNil)

	loginDefsPath = s.tempDir + "/login.defs"
	ioutil.WriteFile(loginDefsPath, []byte(""+
		"CRACKLIB_DICTPATH /usr/lib64/cracklib_dict\n"+
		"\n"+
		"#\n"+
		"# Min/max values for automatic uid selection in useradd\n"+
		"#\n"+
		"UID_MIN			 1000\n"+
		"UID_MAX			60000\n"+
		"# System accounts\n"+
		"SYS_UID_MIN		  101\n"+
		"SYS_UID_MAX		  999\n"+
		"\n"+
		"#\n"+
		"# Min/max values for automatic gid selection in groupadd\n"+
		"#\n"+
		"GID_MIN			 1000\n"), 0644)

	passwdPath = s.tempDir + "/passwd"
	ioutil.WriteFile(passwdPath, []byte(""+
		"root:x:0:0:root:/root:/bin/bash\n"+
		"core:x:500:500:CoreOS Admin:/home/core:/bin/bash\n"+
		"alice:x:1000:500:CoreOS Admin:/home/core:/bin/bash\n"+
		"charlie:x:1001:500:CoreOS Admin:/home/core:/bin/bash\n"+
		""), 0644)

	// replace the real run command with a fake one
	runCommand = func(name string, args ...string) error {
		s.ExecCalls = append(s.ExecCalls, ExecCall{Name: name, Args: args})
		for _, arg := range args {
			if arg == "badluckbrian" {
				return fmt.Errorf("Bad Luck, Brian! Cannot do it")
			}
		}
		return nil
	}
}

func (s *TestUsers) TearDownTest(c *C) {
	os.RemoveAll(s.tempDir)
}

func (s *TestUsers) TestCanSync(c *C) {
	err := SyncUsers(s.DB, s.LoginGroups, false, ioutil.Discard)
	c.Assert(s.ExecCalls, DeepEquals, []ExecCall{
		ExecCall{Name: "useradd", Args: []string{"-c", "", "-p", "*", "--create-home", "badluckbrian"}},
		ExecCall{Name: "useradd", Args: []string{"-c", "", "-p", "*", "--create-home", "dave"}},
		ExecCall{Name: "userdel", Args: []string{"charlie"}},
	})
	c.Assert(err.Error(), Equals, "badluckbrian: Bad Luck, Brian! Cannot do it")
}

func (s *TestUsers) TestBadLoginDefs(c *C) {
	os.Remove(loginDefsPath)
	min, max, err := getUserIDRange()
	c.Assert(min, Equals, int64(1000))
	c.Assert(max, Equals, int64(60000))
	c.Assert(err, IsNil)

	ioutil.WriteFile(loginDefsPath, []byte(""), 0644)
	min, max, err = getUserIDRange()
	c.Assert(min, Equals, int64(-1))
	c.Assert(max, Equals, int64(-1))
	c.Assert(err, ErrorMatches, "Cannot find UID_MIN .*")

	ioutil.WriteFile(loginDefsPath, []byte("UID_MIN 1000\n"), 0644)
	min, max, err = getUserIDRange()
	c.Assert(min, Equals, int64(-1))
	c.Assert(max, Equals, int64(-1))
	c.Assert(err, ErrorMatches, "Cannot find UID_MAX .*")

	ioutil.WriteFile(loginDefsPath, []byte("UID_MIN foo\n"), 0644)
	min, max, err = getUserIDRange()
	c.Assert(min, Equals, int64(-1))
	c.Assert(max, Equals, int64(-1))
	c.Assert(err, ErrorMatches, ".*invalid syntax")

	ioutil.WriteFile(loginDefsPath, []byte("UID_MIN 1000\nUID_MAX foo"), 0644)
	min, max, err = getUserIDRange()
	c.Assert(min, Equals, int64(-1))
	c.Assert(max, Equals, int64(-1))
	c.Assert(err, ErrorMatches, ".*invalid syntax")

	loginDefsPath = "/dev/null/~"
	min, max, err = getUserIDRange()
	c.Assert(min, Equals, int64(-1))
	c.Assert(max, Equals, int64(-1))
	c.Assert(err, ErrorMatches, ".*: not a directory")
}

func (s *TestUsers) TestBadPasswd(c *C) {
	ioutil.WriteFile(passwdPath, []byte(""+
		"root:x:0:0:root:/root:/bin/bash\n"+
		"core:x:500:500:CoreOS Admin:/home/core:/bin/bash\n"+
		"alice:x:1000:500:CoreOS Admin:/home/core:/bin/bash\n"+
		"bob:x:1001:500:CoreOS Admin:/home/core:/bin/bash\n"+
		"charlie:x:1001:500:CoreOS Admin:/home/core:/bin/bash\n"+
		""), 0644)
	users, err := getActualUserNames()
	c.Assert(err, IsNil)
	c.Assert(users, DeepEquals, map[string]struct{}{
		"alice":   struct{}{},
		"bob":     struct{}{},
		"charlie": struct{}{},
	})

	ioutil.WriteFile(passwdPath, []byte(""+
		"root:x:0:0:root:/root:/bin/bash\n"+
		"core:x:500:500:CoreOS Admin:/home/core:/bin/bash\n"+
		"alice:x\n"+
		"bob:x:FOO:500:CoreOS Admin:/home/core:/bin/bash\n"+
		"charlie:x:1001:500:CoreOS Admin:/home/core:/bin/bash\n"+
		""), 0644)
	users, err = getActualUserNames()
	c.Assert(err, ErrorMatches, "invalid line in .*FOO.* invalid syntax.")
	c.Assert(users, DeepEquals, map[string]struct{}{
		"charlie": struct{}{},
	})

	os.Remove(passwdPath)
	users, err = getActualUserNames()
	c.Assert(err, ErrorMatches, ".*no such file or directory")
}
