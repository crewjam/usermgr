package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/crewjam/usermgr"
)

var catCommand = cli.Command{
	Name:   "cat",
	Usage:  "Show an account",
	Action: WithError(CatCommand),
}

func CatCommand(ctx *cli.Context) error {
	config, err := LoadConfig(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	userData, err := usermgr.GetLocalCache(config.CacheDir, config.HostKey)
	if err != nil {
		return err
	}

	name := ctx.Args().First()
	if name != "" {
		user := userData.GetUserByName(name)
		if user == nil {
			return fmt.Errorf("%s: not found\n", name)
		}

		buf, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return err
		}
		ctx.App.Writer.Write(buf)
		return nil
	}

	buf, err := json.MarshalIndent(userData, "", "  ")
	if err != nil {
		return err
	}
	ctx.App.Writer.Write(buf)
	return nil

}
