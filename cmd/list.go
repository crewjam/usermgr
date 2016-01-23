package cmd

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/crewjam/usermgr"
)

var listCommand = cli.Command{
	Name:   "list",
	Usage:  "List all user accounts",
	Action: WithError(ListCommand),
}

func ListCommand(ctx *cli.Context) error {
	config, err := LoadConfig(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	userData, err := usermgr.GetLocalCache(config.CacheDir, config.HostKey)
	if err != nil {
		return err
	}

	for _, user := range userData.Users {
		fmt.Fprintf(ctx.App.Writer, "%s\n", user.Name)
	}
	return nil
}
