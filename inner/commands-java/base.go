package commands_java

import (
	"fmt"
	"github.com/FirewineXie/govm/inner/config"
	web_go "github.com/FirewineXie/govm/inner/web-go"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
)

var configLocal = config.Default().Settings[config.JAVA]

func CommandUninstall(ctx *cli.Context) error {
	versionS := ctx.Args().First()

	version := getCurrentVersion()
	if versionS == version {
		return cli.NewExitError("不能卸载当前版本", 1)
	}
	err := os.RemoveAll(filepath.Join(configLocal.Downloads, "java"+versionS))
	if err != nil {
		return cli.NewExitError("删除该版本失败+"+err.Error(), 1)
	}
	fmt.Println("finish uninstall")
	return nil
}

// CommandUse 激活使用
func CommandUse(ctx *cli.Context) error {
	v, err := getVersion(ctx, true)
	if err != nil {

		return err
	}
	// active use
	_ = os.Remove(configLocal.Symlink)
	fmt.Println(path.Join(configLocal.Downloads, "java"+v), configLocal.Symlink)
	if err := os.Symlink(path.Join(configLocal.Downloads, "java"+v), configLocal.Symlink); err != nil {
		return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
	}
	output, err := exec.Command("java", "version").Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// CommandListRemote 获取远程的可下载的版本
func CommandListRemote(ctx *cli.Context) error {
	versionType := ctx.Args().First()

	collector, err := web_go.NewCollector("")
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

	v := getInstalled(configLocal.Downloads)

	for i := 0; i < len(v); i++ {
		version := v[i]

		str := ""
		goVersion := fmt.Sprintf("go%v", version)
		if in == goVersion {
			str = str + "  * "
		} else {
			str = str + "    "
		}
		str = str + regexp.MustCompile("java").ReplaceAllString(version, "")
		if in == goVersion {
			str = str + " (Currently using " + in + " executable)"
		}
		fmt.Printf(str + "\n")

	}
	if len(v) == 0 {
		fmt.Println("No installations recognized.")
	}
}
