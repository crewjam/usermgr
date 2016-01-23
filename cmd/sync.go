package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/crewjam/usermgr"
)

var SyncInterval = time.Minute * 9

var syncCommand = cli.Command{
	Name:  "sync",
	Usage: "Synchronize local configuration files",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "don't actually change anything, just print what would be changed",
		},
	},
	Action: WithError(Sync),
}

func SyncOnce(config *Config, dryRun bool, stdout io.Writer) error {
	usersData, err := usermgr.UpdateLocalCache(config.CacheDir, config.URL, config.HostKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "updating user data: %s\n", err)
		usersData, err = usermgr.GetLocalCache(config.CacheDir, config.HostKey)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "using cached data\n")
	}
	return syncOnce(config, usersData, dryRun, stdout)
}

func Sync(ctx *cli.Context) error {
	config, err := LoadConfig(ctx.GlobalString("config"))
	if err != nil {
		return err
	}
	// TODO(ross): implement daemon mode
	return SyncOnce(config, ctx.Bool("dry-run"), ctx.App.Writer)
}
