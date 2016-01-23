package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/crewjam/usermgr/cmd"
)

type ErrorWithExitCode interface {
	error
	ExitCode() int
}

func main() {
	// A small hack to deal with cases where the we must be invoked with no
	// arguments, such as when we are invoked as an AuthorizedKeysCommand or
	// a login shell. If the name of the program contains dots, then everything
	// after the first '.' becomes an argument. For for example,
	// `usermgr.sshkeys bob` is equivalent to `usermgr sshkeys bob`. (Of course
	// you need a symlink from usermgr to usermgr.sshkeys to make this work).
	args := strings.Split(os.Args[0], ".")
	args = append(args, os.Args[1:]...)

	if err := cmd.Main(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		if e, ok := err.(ErrorWithExitCode); ok {
			os.Exit(e.ExitCode())
		} else {
			os.Exit(1)
		}
	}
}
