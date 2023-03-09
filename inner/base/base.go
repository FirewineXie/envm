package base

import (
	"fmt"
	"github.com/FirewineXie/govm/inner/arch"
	"github.com/FirewineXie/govm/inner/config"
	"github.com/FirewineXie/govm/inner/web"
	"github.com/mholt/archiver"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
)

func CommandUninstall(ctx *cli.Context) error {
	versionS := ctx.Args().First()

	version := getCurrentVersion()
	if versionS == version {
		return cli.NewExitError("不能卸载当前版本", 1)
	}
	err := os.RemoveAll(filepath.Join(config.Default().Downloads, "go"+versionS))
	if err != nil {
		return cli.NewExitError("删除该版本失败+"+err.Error(), 1)
	}
	fmt.Println("finish uninstall")
	return nil
}

// CommandInstall 安装命令
func CommandInstall(ctx *cli.Context) error {
	versionS := ctx.Args().First()
	collector, err := web.NewCollector("")
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("collect version error1 + %v", err), 1)
	}
	var version *web.Version
	versions, err := collector.AllVersions()
	for _, v := range versions {
		if v.Name == versionS {
			version = v
			break
		}
	}
	findPackage, err := version.FindPackage(web.ArchiveKind, runtime.GOOS, arch.Validate())
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("find version of system error + %v", err), 1)
	}
	downloadPath := filepath.Clean(filepath.Join(config.Default().Downloads, findPackage.FileName))
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
	unchivePath := filepath.Clean(filepath.Join(config.Default().Downloads))
	if err = archiver.Unarchive(downloadPath, unchivePath); err != nil {
		return cli.NewExitError(fmt.Sprintf(" %s", err.Error()), 1)
	}
	err = os.Remove(downloadPath)
	// 目录重命名
	if err = os.Rename(filepath.Join(unchivePath, "go"), filepath.Clean(filepath.Join(config.Default().Downloads, "go"+versionS))); err != nil {
		return cli.NewExitError(fmt.Sprintf(" %s", err.Error()), 1)
	}
	fmt.Println("Installed successfully")
	return nil
}

// CommandUse 激活使用go版本
func CommandUse(ctx *cli.Context) error {
	v, err := getVersion(ctx, true)
	if err != nil {

		return err
	}
	// active use
	_ = os.Remove(config.Default().Symlink)
	fmt.Println(path.Join(config.Default().Downloads, "go"+v), config.Default().Symlink)
	if err := os.Symlink(path.Join(config.Default().Downloads, "go"+v), config.Default().Symlink); err != nil {
		return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
	}
	output, err := exec.Command("go", "version").Output()
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

func CommandClearCache(ctx *cli.Context) error {
	return os.Remove(filepath.Join(config.Default().Root, "meta-page.html"))
}

// CommandListInstalled 展示已经安装的go 版本
func CommandListInstalled(ctx *cli.Context) {
	in := getCurrentVersion()

	v := getInstalled(config.Default().Downloads)

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
