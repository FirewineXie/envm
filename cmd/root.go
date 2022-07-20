package cmd

import (
	"fmt"
	"github.com/FirewineXie/govm/inner/config"
	"github.com/urfave/cli"
	"os"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	app := cli.NewApp()
	app.Name = "govm"
	app.Usage = "Golang Version Manager"
	app.Version = "v1.0.0"

	app.Authors = []cli.Author{
		cli.Author{
			Name: "Firewine",
		},
	}
	// 加载配置
	app.Before = func(ctx *cli.Context) (err error) {

		config.ReadSettings()

		err = os.MkdirAll(config.Default().Download, 0777)
		if err != nil {
			config.SaveSettings()
			return cli.NewExitError("download path error", 1)
		}

		return nil
	}
	app.Commands = commands

	app.After = func(ctx *cli.Context) error {

		config.SaveSettings()
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[g] %s\n", err.Error())
		os.Exit(1)
	}
}
