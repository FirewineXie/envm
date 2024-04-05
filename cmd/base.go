package cmd

import (
	"github.com/FirewineXie/govm/inner/arch"
	"github.com/FirewineXie/govm/inner/commands-go"
	commands_java "github.com/FirewineXie/govm/inner/commands-java"
	"github.com/FirewineXie/govm/inner/config"
	"github.com/urfave/cli"
)

var (
	baseCommands = []cli.Command{
		{
			Name:      "arch",
			Usage:     "systemc arch",
			UsageText: "envm arch",
			Action:    arch.CommandArch,
		},
		{
			Name:      "go",
			Usage:     "env go",
			UsageText: "envm go",
			Before: func(context *cli.Context) error {
				// 校验 go env 是否已经配置
				return config.VerifyEnvGo()

			},
			Subcommands: goCommands,
		},
		{
			Name:      "java",
			Usage:     "env java",
			UsageText: "envm java",
			Before: func(context *cli.Context) error {
				// 校验 java env 是否已经配置
				return config.VerifyEnvJava()
			},
			Subcommands: javaCommands,
		},
	}

	goCommands = []cli.Command{
		{
			Name:      "ls",
			Usage:     "List installed versions",
			UsageText: "govm ls",
			Action:    commands_go.CommandListInstalled,
		},
		{
			Name:      "lsr",
			Usage:     "List remote versions available for install",
			UsageText: "govm lsr [stable|archived]",
			Action:    commands_go.CommandListRemote,
		},
		{
			Name:      "active",
			Usage:     "Switch to specified version",
			UsageText: "govm active <version>",
			Action:    commands_go.CommandUse,
		},
		{
			Name:      "install",
			Usage:     "Download and install a <version>",
			UsageText: "govm install <version>",
			Action:    commands_go.CommandInstall,
		},
		{
			Name:      "uninstall",
			Usage:     "Uninstall a version",
			UsageText: "gvm uninstall <version>",
			Action:    commands_go.CommandUninstall,
		},
	}

	javaCommands = []cli.Command{
		{
			Name:      "ls",
			Usage:     "List installed versions",
			UsageText: "envm java  ls",
			Action:    commands_java.CommandListInstalled,
		},
		{
			Name:      "active",
			Usage:     "Switch to specified version",
			UsageText: "envm java active <version>",
			Action:    commands_java.CommandUse,
		},
		{
			Name:      "uninstall",
			Usage:     "Uninstall a version",
			UsageText: "envm java uninstall <version>",
			Action:    commands_java.CommandUninstall,
		},
	}
)
