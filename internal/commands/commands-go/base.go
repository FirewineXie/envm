package commands_go

import (
	"fmt"
	"github.com/FirewineXie/envm/internal/arch"
	"github.com/FirewineXie/envm/internal/commands/common"
	"github.com/FirewineXie/envm/internal/config"
	"github.com/FirewineXie/envm/internal/logic/web-go"
	"github.com/FirewineXie/envm/util"
	"github.com/mholt/archiver/v3"

	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var configLocal = config.Default().LinkSetting[config.GO]

func CommandUninstall(ctx *cli.Context) error {
	versionS := ctx.Args().First()

	version := common.GetCurrentVersion("go")
	if versionS == version {
		return cli.NewExitError("不能卸载当前版本", 1)
	}
	err := os.RemoveAll(filepath.Join(configLocal.Downloads, "go"+versionS))
	if err != nil {
		return cli.NewExitError("删除该版本失败+"+err.Error(), 1)
	}
	fmt.Println("finish uninstall")
	return nil
}

// CommandInstall 安装命令
func CommandInstall(ctx *cli.Context) error {
	versionS := ctx.Args().First()
	
	// 检查版本是否已存在
	targetDir := filepath.Clean(filepath.Join(configLocal.Downloads, "go"+versionS))
	if exists, err := util.PathExists(targetDir); err != nil {
		return cli.NewExitError(fmt.Sprintf("check existing version error: %v", err), 1)
	} else if exists {
		fmt.Printf("Go version %s is already installed\n", versionS)
		return nil
	}
	
	collector, err := web_go.NewCollector("")
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("collect version error1 + %v", err), 1)
	}
	var version *web_go.VersionGO
	versions, err := collector.AllVersions()
	for _, v := range versions {
		if v.Name == versionS {
			version = v
			break
		}
	}
	findPackage, err := version.FindPackage(util.ArchiveKind, runtime.GOOS, arch.Validate())
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("find version of system error + %v", err), 1)
	}
	downloadPath := filepath.Clean(filepath.Join(configLocal.Downloads, findPackage.FileName))
	
	// 检查下载文件是否已存在，避免重复下载
	if exists, err := util.PathExists(downloadPath); err != nil {
		return cli.NewExitError(fmt.Sprintf("check download file error: %v", err), 1)
	} else if !exists {
		// 根据收集器的URL确定下载前缀
		downloadPrefix := "https://golang.org"
		if strings.Contains(collector.GetURL(), "golang.google.cn") {
			downloadPrefix = "https://golang.google.cn"
		}
		findPackage.URL = downloadPrefix + findPackage.URL
		fmt.Printf("Downloading Go %s...\n", versionS)
		err = findPackage.DownloadV2(downloadPath)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("download version error + %v", err), 1)
		}
		err = findPackage.VerifyChecksum(downloadPath)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("verify version error + %v", err), 1)
		}
	} else {
		fmt.Printf("Archive file already exists, skipping download\n")
	}

	// 解压安装包
	unchivePath := filepath.Clean(filepath.Join(configLocal.Downloads))
	if err = archiver.Unarchive(downloadPath, unchivePath); err != nil {
		return cli.NewExitError(fmt.Sprintf(" %s", err.Error()), 1)
	}
	err = os.Remove(downloadPath)
	// 目录重命名
	if err = os.Rename(filepath.Join(unchivePath, "go"), targetDir); err != nil {
		return cli.NewExitError(fmt.Sprintf(" %s", err.Error()), 1)
	}
	fmt.Println("Installed successfully")
	return nil
}

// CommandUse 激活使用go版本
func CommandUse(ctx *cli.Context) error {
	v, err := common.GetVersion(ctx, configLocal.Downloads, "go", true)
	if err != nil {
		return err
	}
	// active use
	targetPath := path.Join(configLocal.Downloads, "go"+v)
	fmt.Printf("Switching to Go version %s\n", v)
	fmt.Printf("Creating symlink: %s -> %s\n", configLocal.Symlink, targetPath)
	
	if err := util.CreateSymlink(targetPath, configLocal.Symlink); err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to create symlink: %s", err.Error()), 1)
	}
	
	output, err := exec.Command("go", "version").Output()
	if err != nil {
		return err
	}
	fmt.Printf("Successfully switched to: %s", string(output))
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
func CommandListInstalled(ctx *cli.Context) error {
	in := common.GetCurrentVersion("go")

	v := common.GetInstalled(configLocal.Downloads, "go")

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
		fmt.Print(str + "\n")

	}
	if len(v) == 0 {
		fmt.Println("No installations recognized.")
	}
	return nil
}
