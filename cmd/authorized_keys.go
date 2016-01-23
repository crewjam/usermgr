package cmd

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/crewjam/usermgr"
)

var authorizedKeysCommand = cli.Command{
	Name:   "authorized-keys",
	Usage:  "Print the SSH authorized keys for a user",
	Action: WithError(AuthorizedKeys),
}

// AuthorizedKeys implements the "authorized-keys" subcommand which prints the
// authorized keys for the specified user to stdout.
//
// From the sshd_config(8) manpage:
//
//     AuthorizedKeysCommand
//             Specifies a program to be used to look up the user's public keys.  The program
//             will be invoked with a single argument of the username being authenticated, and
//             should produce on standard output zero or more lines of authorized_keys output
//             (see AUTHORIZED_KEYS in sshd(8)).  If a key supplied by AuthorizedKeysCommand does
//             not successfully authenticate and authorize the user then public key authentica-
//             tion continues using the usual AuthorizedKeysFile files.  By default, no Autho-
//             rizedKeysCommand is run.
//
//     AuthorizedKeysCommandUser
//             Specifies the user under whose account the AuthorizedKeysCommand is run.  It is
//             recommended to use a dedicated user that has no other role on the host than run-
//             ning authorized keys commands.
//
// Example configuration:
//
//  AuthorizedKeysCommand /usr/bin/usermgr.sshkeys
//  AuthorizedKeysCommandUser nobody
//
func AuthorizedKeys(ctx *cli.Context) error {
	config, err := LoadConfig(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	userData, err := usermgr.GetLocalCache(config.CacheDir, config.HostKey)
	if err != nil {
		return err
	}

	name := ctx.Args().First()
	user := userData.GetUserByName(name)
	if user == nil {
		return fmt.Errorf("%s: not found\n", name)
	}

	if user.AuthorizedKeys != nil {
		for _, authorizedKey := range user.AuthorizedKeys {
			fmt.Fprintf(ctx.App.Writer, "%s\n", authorizedKey)
		}
	}
	return nil
}
