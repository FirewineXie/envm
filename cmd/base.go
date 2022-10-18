package cmd

import (
	"github.com/FirewineXie/govm/inner/arch"
	"github.com/FirewineXie/govm/inner/base"
	"github.com/FirewineXie/govm/inner/config"
	"github.com/urfave/cli"
)

var (
	commands = []cli.Command{
		{
			Name:        "config",
			Usage:       "config downloads",
			UsageText:   "govm config downloads",
			Description: "only set ,not judge it is valid",
			Action:      config.SetDownloads,
		},
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
			Name:   "clear",
			Usage:  "clear meta-data html",
			Action: base.CommandClearCache,
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
