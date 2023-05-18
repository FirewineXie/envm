package commands_java

import (
	"fmt"
	"github.com/FirewineXie/govm/inner/arch"
	"github.com/FirewineXie/govm/inner/config"
	web_java "github.com/FirewineXie/govm/inner/web-java"
	"github.com/mholt/archiver"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
)

var configLocal = config.Default().Settings[config.JAVA]

func CommandUninstall(ctx *cli.Context) error {
	versionS := ctx.Args().First()

	version := getCurrentVersion()
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
	v, err := getVersion(ctx, true)
	if err != nil {
		return err
	}
	// active use
	_ = os.Remove(configLocal.Symlink)
	fmt.Println(path.Join(configLocal.Downloads, "jdk-"+v), configLocal.Symlink)
	if err := os.Symlink(path.Join(configLocal.Downloads, "jdk-"+v), configLocal.Symlink); err != nil {
		return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
	}
	output, err := exec.Command("java", "version").Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// CommandListInstalled 展示已经安装
func CommandListInstalled(ctx *cli.Context) {
	in := getCurrentVersion()

	v := getInstalled(configLocal.Downloads)

	for i := 0; i < len(v); i++ {
		version := v[i]

		str := ""
		goVersion := fmt.Sprintf("java%v", version)
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

func CommandListRemote(ctx *cli.Context) error {

	collector, err := web_java.NewCollector("")
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("collect version error1 + %v", err), 1)
	}
	items, err := collector.LatestFiveVersion()
	if err != nil {
		return err
	} else {
		for _, version := range items {
			fmt.Println(version.Name)
		}

		return nil
	}

}

func CommandInstall(ctx *cli.Context) error {
	versionS := "20"
	collector, err := web_java.NewCollector("")
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("collect version error1 + %v", err), 1)
	}
	var version *web_java.Version
	versions, err := collector.LatestFiveVersion()
	for _, v := range versions {
		if v.Name == versionS {
			version = v
			break
		}
	}
	findPackage, err := version.FindPackage(versionS, runtime.GOOS, arch.Validate())
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("find version of system error + %v", err), 1)
	}
	downloadPath := filepath.Clean(filepath.Join(configLocal.Downloads, findPackage.FileName))
	findPackage.URL = "https://golang.google.cn" + findPackage.URL
	err = findPackage.DownloadV2(downloadPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("download version error + %v", err), 1)
	}
	err = findPackage.VerifyChecksum(downloadPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("verify version error + %v", err), 1)
	}

	// 解压安装包
	unchivePath := filepath.Clean(filepath.Join(configLocal.Downloads))
	if err = archiver.Unarchive(downloadPath, unchivePath); err != nil {
		return cli.NewExitError(fmt.Sprintf(" %s", err.Error()), 1)
	}
	err = os.Remove(downloadPath)
	// 目录重命名
	if err = os.Rename(filepath.Join(unchivePath, "go"), filepath.Clean(filepath.Join(configLocal.Downloads, "go"+versionS))); err != nil {
		return cli.NewExitError(fmt.Sprintf(" %s", err.Error()), 1)
	}
	fmt.Println("Installed successfully")
	return nil
}
