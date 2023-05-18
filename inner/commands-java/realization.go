package commands_java

import (
	"errors"
	"github.com/blang/semver"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

// 获取版本
func getVersion(ctx *cli.Context, localInstallsOnly ...bool) (string, error) {
	version := ctx.Args().First()
	if version == "" {
		return "", cli.ShowSubcommandHelp(ctx)
	}

	// 如果是true ,那么本地寻找 这个版本
	if localInstallsOnly[0] {
		installed := getInstalled(configLocal.Downloads)
		for _, installVersion := range installed {
			if installVersion == version {
				return version, nil
			}
		}
	}
	return "", errors.New("you have not install it,please install before use")
}

// 获取当前版本
func getCurrentVersion() (version string) {
	cmd := exec.Command("java", "--version")
	str, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	split := strings.Split(string(str), "\n")
	if len(split) > 1 {
		itemp := strings.Split(split[0], " ")
		return itemp[1]
	}
	return string(str)

}

// 获取下载的版本列表
func getInstalled(root string) []string {
	list := make([]semver.Version, 0)
	entries, _ := os.ReadDir(path.Clean(root))

	for _, entry := range entries {

		if entry.IsDir() {
			isGo, _ := regexp.MatchString("jdk", entry.Name())
			if isGo {
				currentVersionString := strings.Replace(entry.Name(), "jdk-", "", 1)
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
