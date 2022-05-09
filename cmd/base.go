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
			Name:      "ls-remote",
			Usage:     "List remote versions available for install",
			UsageText: "gvm ls-remote [stable|archived]",
			Action:    base.CommandListRemote,
		},
		{
			Name:      "active",
			Usage:     "Switch to specified version",
			UsageText: "gvm active <version>",
			Action:    base.CommandUse,
		},
		{
			Name:      "install",
			Usage:     "Download and install a <version>",
			UsageText: "gvm install <version>",
			Action:    install,
		},
		{
			Name:      "uninstall",
			Usage:     "Uninstall a version",
			UsageText: "gvm uninstall <version>",
			Action:    uninstall,
		},
	}
)
