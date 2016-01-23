package usermgr

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/crewjam/errset"
)

var loginDefsPath = "/etc/login.defs"
var passwdPath = "/etc/passwd"

var runCommand = func(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getUserIDRange() (min int64, max int64, err error) {
	min, max = -1, -1
	f, err := os.Open(loginDefsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Default values of UID_MIN and UID_MAX (on linux)
			// in the case that login.defs is missing
			return 1000, 60000, nil
		}
		return -1, -1, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	re := regexp.MustCompile(`^\s*(\w+)\s+(\w+)\s*$`)
	for s.Scan() {
		m := re.FindStringSubmatch(s.Text())
		if m == nil {
			continue
		}
		name, value := m[1], m[2]

		if name == "UID_MAX" {
			max, err = strconv.ParseInt(value, 10, 32)
			if err != nil {
				return -1, -1, err
			}
		}

		if name == "UID_MIN" {
			min, err = strconv.ParseInt(value, 10, 32)
			if err != nil {
				return -1, -1, err
			}
		}

	}
	if err := s.Err(); err != nil {
		return -1, -1, err
	}
	if min == -1 {
		return -1, -1, fmt.Errorf("Cannot find UID_MIN in %s", loginDefsPath)
	}
	if max == -1 {
		return -1, -1, fmt.Errorf("Cannot find UID_MAX in %s", loginDefsPath)
	}
	return min, max, err
}

func getActualUserNames() (map[string]struct{}, error) {
	min, max, err := getUserIDRange()
	if err != nil {
		return nil, err
	}

	rv := map[string]struct{}{}
	f, err := os.Open(passwdPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	errs := errset.ErrSet{}
	for s.Scan() {
		parts := strings.Split(s.Text(), ":")
		if len(parts) < 3 {
			errs = append(errs, fmt.Errorf("invalid line in %s: %q", passwdPath, s.Text()))
			continue
		}
		userName := parts[0]
		uid, err := strconv.ParseInt(parts[2], 10, 32)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid line in %s: %q (%s)", passwdPath, s.Text(), err))
			continue
		}
		if uid >= min && uid <= max {
			rv[userName] = struct{}{}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return rv, errs.ReturnValue()
}

func SyncUsers(ud *UsersData, groups []string, dryRun bool, stdout io.Writer) error {
	errs := errset.ErrSet{}

	actualUserNames, err := getActualUserNames()
	if err != nil {
		errs = append(errs, fmt.Errorf("list users: %s", err))
	}
	if actualUserNames == nil {
		actualUserNames = map[string]struct{}{}
	}

	for _, nominalUser := range ud.Users {
		if !nominalUser.InAnyGroup(groups) {
			continue
		}

		if _, userExists := actualUserNames[nominalUser.Name]; userExists {
			continue
		}

		fmt.Fprintf(stdout, "%s: create\n", nominalUser.Name)
		if dryRun {
			continue
		}
		err := runCommand("useradd", "-c", nominalUser.RealName,
			"-p", "*", "--create-home", nominalUser.Name)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %s", nominalUser.Name, err))
		}
	}

	for actualUserName := range actualUserNames {
		existingNominalUser := ud.GetUserByName(actualUserName)
		if existingNominalUser != nil && existingNominalUser.InAnyGroup(groups) {
			continue
		}

		fmt.Fprintf(stdout, "%s: remove\n", actualUserName)
		if dryRun {
			continue
		}
		if err := runCommand("userdel", actualUserName); err != nil {
			errs = append(errs, fmt.Errorf("%s: %s", actualUserName, err))
		}
	}

	return errs.ReturnValue()
}
