package cmd

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/crewjam/usermgr"
)

var keygenCommand = cli.Command{
	Name:   "keygen",
	Usage:  "Generate a keypair",
	Action: WithError(KeygenCommand),
}

func KeygenCommand(ctx *cli.Context) error {
	adminKey := usermgr.GenerateKeyPair()
	fmt.Fprintf(ctx.App.Writer, "admin key: %s\n", adminKey)
	fmt.Fprintf(ctx.App.Writer, "host key: %s\n", adminKey.HostKey)
	return nil
}
