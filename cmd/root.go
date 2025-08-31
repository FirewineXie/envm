package cmd

import (
	"fmt"
	"github.com/FirewineXie/envm/internal/config"
	"github.com/FirewineXie/envm/util"

	"github.com/urfave/cli"
	"os"
)

// Execute adds all child goCommands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// 检查是否是管理员权限命令（Windows UAC提升后的执行）
	if len(os.Args) > 1 && util.IsAdminSymlinkCommand(os.Args[1:]) {
		err := util.HandleAdminSymlinkCommand(os.Args[1:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "管理员权限命令执行失败: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println("管理员权限命令执行成功")
		return
	}

	app := cli.NewApp()
	app.Name = "envm"
	app.Usage = "Any More Version Manager"
	app.Version = "v1.0.2"
	app.Description = `
			java & go version manager
     `

	app.Authors = []cli.Author{
		cli.Author{
			Name: "Firewine",
		},
	}
	app.Before = func(context *cli.Context) error {
		return config.VerifyEnv()
	}

	app.Commands = baseCommands

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[g] %s\n", err.Error())
		os.Exit(1)
	}
}
