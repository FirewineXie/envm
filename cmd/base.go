package cmd

import (
	"github.com/FirewineXie/govm/inner/arch"
	"github.com/FirewineXie/govm/inner/base"
	"github.com/urfave/cli"
)

var (
	commands = []cli.Command{

		{
			Name:      "arch",
			Usage:     "systemc arch",
			UsageText: "govm arch",
			Action:    arch.CommandArch,
		},
		{
			Name:      "ls",
			Usage:     "List installed versions",
			UsageText: "govm ls",
			Action:    base.CommandListInstalled,
		},
		{
			Name:      "lsr",
			Usage:     "List remote versions available for install",
			UsageText: "govm lsr [stable|archived]",
			Action:    base.CommandListRemote,
		},
		{
			Name:      "active",
			Usage:     "Switch to specified version",
			UsageText: "govm active <version>",
			Action:    base.CommandUse,
		},
		{
			Name:      "install",
			Usage:     "Download and install a <version>",
			UsageText: "govm install <version>",
			Action:    base.CommandInstall,
		},
		{
			Name:      "uninstall",
			Usage:     "Uninstall a version",
			UsageText: "gvm uninstall <version>",
			Action:    base.CommandUninstall,
		},
	}
)
