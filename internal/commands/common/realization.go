package common

import (
	"errors"
	"github.com/blang/semver/v4"
	"github.com/urfave/cli"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

// 获取版本
func GetVersion(ctx *cli.Context, downloadPath string, localInstallsOnly ...bool) (string, error) {
	version := ctx.Args().First()
	if version == "" {
		return "", cli.ShowSubcommandHelp(ctx)
	}

	// 如果是true ,那么本地寻找 这个版本
	if localInstallsOnly[0] {
		installed := GetInstalled(downloadPath)
		for _, installVersion := range installed {
			if installVersion == version {
				return version, nil
			}
		}
	}
	return "", errors.New("you have not install it,please install before use")
}

// 获取当前版本
func GetCurrentVersion() (version string) {
	cmd := exec.Command("node", "version")
	str, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	split := strings.Split(string(str), " ")
	if len(split) > 3 {
		return split[2]
	}
	return string(str)

}

// 获取下载的版本列表
func GetInstalled(root string) []string {
	list := make([]semver.Version, 0)
	files, _ := ioutil.ReadDir(path.Clean(root))
	for i := len(files) - 1; i >= 0; i-- {
		if files[i].IsDir() {
			isGo, _ := regexp.MatchString("node", files[i].Name())
			if isGo {
				currentVersionString := strings.Replace(files[i].Name(), "node", "", 1)
				if currentVersion, err := semver.Make(currentVersionString); err == nil {
					list = append(list, currentVersion)
				}

			}
		}
	}

	semver.Sort(list)

	loggableList := make([]string, 0)

	for _, version := range list {
		loggableList = append(loggableList, version.String())
	}
	loggableList = reverseStringArray(loggableList)
	return loggableList
}

func reverseStringArray(str []string) []string {
	for i := 0; i < len(str)/2; i++ {
		j := len(str) - i - 1
		str[i], str[j] = str[j], str[i]
	}

	return str
}
