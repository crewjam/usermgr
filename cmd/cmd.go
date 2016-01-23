package cmd

import (
	"io"

	"github.com/codegangsta/cli"
)

var globalErr error

func WithError(f func(context *cli.Context) error) func(*cli.Context) {
	return func(ctx *cli.Context) {
		globalErr = f(ctx)
	}
}

func Main(args []string, output io.Writer) error {
	app := cli.NewApp()
	app.Writer = output
	app.Name = "usermgr"
	app.Usage = "manage user accounts"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: DefaultConfigPath,
			Usage: "The path to the configuration file",
		},
	}
	app.Commands = []cli.Command{
		syncCommand,
		authorizedKeysCommand,
		catCommand,
		listCommand,
		shellCommand,
		keygenCommand,
		webCommand,
	}

	globalErr = nil
	if err := app.Run(args); err != nil {
		return err
	}
	return globalErr
}
