package commands_node

import (
	"fmt"
	"github.com/FirewineXie/envm/internal/arch"
	"github.com/FirewineXie/envm/internal/commands/common"
	"github.com/FirewineXie/envm/internal/config"
	"github.com/FirewineXie/envm/internal/logic/web-node"
	"github.com/FirewineXie/envm/util"

	"github.com/mholt/archiver/v3"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
)

var configLocal = config.Default().LinkSetting[config.NODE]

func CommandUninstall(ctx *cli.Context) error {
	versionS := ctx.Args().First()

	version := common.GetCurrentVersion()
	if versionS == version {
		return cli.NewExitError("不能卸载当前版本", 1)
	}
	err := os.RemoveAll(filepath.Join(configLocal.Downloads, "node"+versionS))
	if err != nil {
		return cli.NewExitError("删除该版本失败+"+err.Error(), 1)
	}
	fmt.Println("finish uninstall")
	return nil
}

// CommandInstall 安装命令
func CommandInstall(ctx *cli.Context) error {
	versionS := ctx.Args().First()
	return commandInstall(versionS)
}
func commandInstall(versionS string) error {
	if versionS == "" {
		return cli.NewExitError(fmt.Sprintf("find version for not empty"), 1)
	}
	_, _, _, _, _, _, err := web_node.GetAvailable()
	if err != nil {
		return cli.NewExitError("get mirror version failed"+err.Error(), 1)
	}
	// 1. 验证版本号，是否正确
	element, ok := web_node.GetMeta()[versionS]
	if !ok {
		return cli.NewExitError(fmt.Sprintf("find version for not found"), 1)
	}

	// 3. 此版本是否已经下载，如果已经下载，则忽略
	if getInstalled(versionS) {
		fmt.Println("this version is downloaded")
		return nil
	}

	// 4. 此版本是否有该系统架构当前的版本
	findPackage, err := element.FindPackage(util.ArchiveKind, runtime.GOOS, arch.Validate())
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("find version of system error + %v", err), 1)
	}

	downloadPath := filepath.Clean(filepath.Join(configLocal.Downloads, findPackage.ArchiveName))
	err = findPackage.DownloadV2(downloadPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("download version error + %v", err), 1)
	}
	//err = findPackage.VerifyChecksum(downloadPath)
	//if err != nil {
	//	return cli.NewExitError(fmt.Sprintf("verify version error + %v", err), 1)
	//}

	// 解压安装包
	unchivePath := filepath.Clean(filepath.Join(configLocal.Downloads))
	if err = archiver.Unarchive(downloadPath, unchivePath); err != nil {
		return cli.NewExitError(fmt.Sprintf(" %s", err.Error()), 1)
	}
	err = os.Remove(downloadPath)
	// 目录重命名
	if err = os.Rename(filepath.Join(unchivePath, findPackage.FileName), filepath.Clean(filepath.Join(unchivePath, "node"+versionS))); err != nil {
		return cli.NewExitError(fmt.Sprintf(" %s", err.Error()), 1)
	}
	fmt.Println("Installed successfully")
	return nil
}

// CommandUse 激活使用
func CommandUse(ctx *cli.Context) error {
	v, err := common.GetVersion(ctx, configLocal.Downloads, true)
	if err != nil {

		return err
	}
	// active use
	if configLocal.Symlink == "" {
		return cli.NewExitError("not config symlink", 1)
	}
	_ = os.Remove(configLocal.Symlink)
	fmt.Println(path.Join(configLocal.Downloads, "node"+v), configLocal.Symlink)
	if err := os.Symlink(path.Join(configLocal.Downloads, "node"+v), configLocal.Symlink); err != nil {
		return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
	}
	output, err := exec.Command("node", "version").Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// CommandListRemote 获取远程的可下载的版本
func CommandListRemote(ctx *cli.Context) error {
	versionType := ctx.Args().First()

	all, lts, current, stable, unstable, _, err := web_node.GetAvailable()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
	}
	releases := 20

	switch versionType {
	case "all":
		for i, version := range all {
			if i == releases {
				break
			}
			fmt.Println(version)
		}
	case "lts":
		for i, version := range lts {
			if i == releases {
				break
			}
			fmt.Println(version)
		}
	case "current":
		for i, version := range current {
			if i == releases {
				break
			}
			fmt.Println(version)
		}
	case "stable":
		for i, version := range stable {
			if i == releases {
				break
			}
			fmt.Println(version)
		}
	case "unstable":
		for i, version := range unstable {
			if i == releases {
				break
			}
			fmt.Println(version)
		}

	}

	return cli.ShowSubcommandHelp(ctx)
}

// CommandListInstalled 展示已经安装
func CommandListInstalled(ctx *cli.Context) {
	in := common.GetCurrentVersion()

	v := common.GetInstalled(configLocal.Downloads)

	for i := 0; i < len(v); i++ {
		version := v[i]

		str := ""
		goVersion := fmt.Sprintf("node%v", version)
		if in == goVersion {
			str = str + "  * "
		} else {
			str = str + "    "
		}
		str = str + regexp.MustCompile("node").ReplaceAllString(version, "")
		if in == goVersion {
			str = str + " (Currently using " + in + " executable)"
		}
		fmt.Printf(str + "\n")

	}
	if len(v) == 0 {
		fmt.Println("No installations recognized.")
	}
}

// getInstalled 展示已经安装
func getInstalled(currentVersion string) bool {

	v := common.GetInstalled(configLocal.Downloads)
	for i := 0; i < len(v); i++ {
		version := v[i]
		if version == currentVersion {
			return true
		}
	}
	return false
}
