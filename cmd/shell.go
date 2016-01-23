package cmd

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/crewjam/usermgr"
)

var syscallExec = syscall.Exec
var stdin io.Reader = os.Stdin
var osGetenv = os.Getenv

var shellCommand = cli.Command{
	Name:   "shell",
	Usage:  "Be a login shell",
	Action: WithError(ShellCommand),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "command, c",
			Value: "",
			Usage: "The command to run",
		},
	},
}

func ShellCommand(ctx *cli.Context) error {
	config, err := LoadConfig(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	userData, err := usermgr.GetLocalCache(config.CacheDir, config.HostKey)
	if err != nil {
		return err
	}

	// Print a warning message
	if len(os.Args) > 0 && len(os.Args[0]) > 0 && os.Args[0][0] == '-' {
		fmt.Fprintf(ctx.App.Writer, "[shell logging active]\n")
	}

	user := userData.GetUserByName(os.Getenv("USER"))
	if user == nil {
		return fmt.Errorf("unknown user")
	}

	if config.LoginMFARequried {
		err := func() error {
			for try := 0; try < 3; try++ {
				fmt.Fprintf(ctx.App.Writer, "Multi-factor code: ")
				var code string
				fmt.Fscanln(stdin, &code)

				err := usermgr.ValidateCode(*user, code, userData.YubikeyClientID, userData.YubikeyClientSecret)
				if err == nil {
					return nil
				}
				fmt.Fprintf(ctx.App.Writer, "Incorrect token: %s\n", err)
			}
			return fmt.Errorf("Incorrect token")
		}()
		if err != nil {
			return err
		}
	}

	os.Setenv("SHELL", "/bin/bash")
	args := []string{"sudo", "-u", os.Getenv("USER"), os.Getenv("SHELL")}

	if ctx.IsSet("command") {
		args = append(args, "-c", ctx.String("command"))
	}

	return syscallExec("/usr/bin/sudo", args, nil)
}
