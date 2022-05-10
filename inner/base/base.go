package base

import (
	"fmt"
	"github.com/FirewineXie/govm/inner/config"
	"github.com/FirewineXie/govm/inner/web"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
)

// CommandUse 激活使用go版本
func CommandUse(ctx *cli.Context) error {
	v, err := getVersion(ctx, true)
	if err != nil {

		return err
	}
	// active use
	_ = os.Remove(config.Default().Symlink)

	if err := os.Symlink(path.Join(config.Default().Download, v), config.Default().Symlink); err != nil {
		return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
	}
	output, err := exec.Command(filepath.Join(config.Default().Root, "bin", "go"), "version").Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// CommandListRemote 获取远程的可下载的版本
func CommandListRemote(ctx *cli.Context) error {
	versionType := ctx.Args().First()

	collector, err := web.NewCollector("")
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("collect version error1 + %v", err), 1)
	}
	if versionType == "stable" {
		versions, err := collector.StableVersions()
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("collect version error2 + %v", err), 1)
		}
		for _, version := range versions {
			fmt.Println(version.Name)
		}
		return nil
	}
	if versionType == "archived" {
		versions, err := collector.ArchivedVersions()
		if err != nil {

			return err
		}
		for _, version := range versions {
			fmt.Println(version.Name)
		}
		return nil
	}

	return cli.ShowSubcommandHelp(ctx)

}

// CommandListInstalled 展示已经安装的go 版本
func CommandListInstalled(ctx *cli.Context) {
	in := getCurrentVersion()

	v := getInstalled(config.Default().Download)

	for i := 0; i < len(v); i++ {
		version := v[i]

		str := ""
		goVersion := fmt.Sprintf("go%v", version)
		if in == goVersion {
			str = str + "  * "
		} else {
			str = str + "    "
		}
		str = str + regexp.MustCompile("go").ReplaceAllString(version, "")
		if in == goVersion {
			str = str + " (Currently using " + in + " executable)"
		}
		fmt.Printf(str + "\n")

	}
	if len(v) == 0 {
		fmt.Println("No installations recognized.")
	}
}
