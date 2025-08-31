package commands_java

import (
	"fmt"
	"github.com/FirewineXie/envm/internal/commands/common"
	"github.com/FirewineXie/envm/internal/config"
	"github.com/FirewineXie/envm/internal/logic/web-java"
	"github.com/FirewineXie/envm/util"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
)

var configLocal = config.Default().LinkSetting[config.JAVA]

func CommandUninstall(ctx *cli.Context) error {
	versionS := ctx.Args().First()

	version := common.GetCurrentVersion("java")
	if versionS == version {
		return cli.NewExitError("不能卸载当前版本", 1)
	}
	err := os.RemoveAll(filepath.Join(configLocal.Downloads, "jdk-"+versionS))
	if err != nil {
		return cli.NewExitError("删除该版本失败+"+err.Error(), 1)
	}
	fmt.Println("finish uninstall")
	return nil
}

// CommandUse 激活使用
func CommandUse(ctx *cli.Context) error {
	v, err := common.GetVersion(ctx, configLocal.Downloads, "jdk", true)
	if err != nil {
		return err
	}
	// active use
	targetPath := path.Join(configLocal.Downloads, v)
	fmt.Printf("Switching to Java version %s\n", v)
	fmt.Printf("Creating symlink: %s -> %s\n", configLocal.Symlink, targetPath)
	
	if err := util.CreateSymlink(targetPath, configLocal.Symlink); err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to create symlink: %s", err.Error()), 1)
	}
	
	output, err := exec.Command("java", "--version").Output()
	if err != nil {
		return err
	}
	fmt.Printf("Successfully switched to: %s", string(output))
	return nil
}

// CommandListInstalled 展示已经安装
func CommandListInstalled(ctx *cli.Context) error {
	in := common.GetCurrentVersion("java")

	v := common.GetInstalled(configLocal.Downloads, "jdk")

	for i := 0; i < len(v); i++ {
		version := v[i]

		str := ""
		goVersion := fmt.Sprintf("java%v", version)
		if in == goVersion {
			str = str + "  * "
		} else {
			str = str + "    "
		}
		str = str + regexp.MustCompile("jdk").ReplaceAllString(version, "")
		if in == goVersion {
			str = str + " (Currently using " + in + " executable)"
		}
		fmt.Print(str + "\n")

	}
	if len(v) == 0 {
		fmt.Println("No installations recognized.")
	}
	return nil
}

func CommandListRemote(ctx *cli.Context) error {

	collector, err := web_java.NewCollector("")
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("collect version error1 + %v", err), 1)
	}
	items, err := collector.LatestFiveVersion()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("collect version error1 + %v", err), 1)
	} else {
		for i, version := range items {
			if i == 20 {
				break
			}
			fmt.Println(version.Name)
		}

	}
	fmt.Println("detail see website")
	return nil
}

func CommandInstall(ctx *cli.Context) error {

	fmt.Println("Installed successfully")
	return nil
}
