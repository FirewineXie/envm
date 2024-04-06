package web_node

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FirewineXie/envm/util"
	"strings"
)

const (
	// DefaultURL 默认镜像地址
	DefaultURL = "https://nodejs.org/dist/"
)

var meta map[string]VersionNode

func GetMeta() map[string]VersionNode {
	return meta
}

type VersionNode struct {
	util.Version
}

// FindPackage 返回指定操作系统和硬件架构的版本包
func (v *VersionNode) FindPackage(kind, goos, goarch string) (*util.Package, error) {
	goarch = exchangeArch(goarch)
	for _, packageData := range v.Packages {
		if packageData.OS == goos && packageData.Kind == kind && goarch == packageData.Arch {
			return packageData, nil
		}
	}

	return nil, util.ErrPackageNotFound
}

// GetAvailable Retrieve the remotely available versions
func GetAvailable() (all []string, lts []string, current []string, stable []string, unstable []string, npm map[string]string, err error) {
	meta = make(map[string]VersionNode)
	resp, errr := DownloadContent(DefaultURL + "index.json")
	if errr != nil {
		err = errors.New("getting mirrors " + errr.Error())
		return
	}
	// Check the service to make sure the version is available
	if len(resp) == 0 {
		err = errors.New("retrieving version list: \"" + DefaultURL + "index.json" + "\" returned blank results. This can happen when the remote file is being updated. Please try again in a few minutes")
		return
	}

	// Parse
	var data []FileData
	errr = json.Unmarshal(resp, &data)
	if errr != nil {
		err = errors.New("retrieving version " + errr.Error())
		return
	}
	npm = make(map[string]string, len(data))
	for _, element := range data {

		var version = element.Version[1:]
		all = append(all, version)
		{
			nodeVersion := VersionNode{
				Version: util.Version{
					Name: version,
				},
			}

			//https://nodejs.org/dist/v20.12.1/node-v20.12.1-darwin-arm64.tar.gz
			for _, file := range element.Files {
				typeFile := ""
				split := strings.Split(file, "-")

				switch split[0] {
				case "darwin":
				case "linux":
					typeFile = "tar.gz"
				case "win":
					typeFile = "zip"
				}
				if typeFile == "" {
					continue
				}
				nodeVersion.Packages = append(nodeVersion.Packages, &util.Package{
					ArchiveName: "node" + version + "." + typeFile,
					FileName: fmt.Sprintf("node-v%s-%s", version, func() (result string) {

						if strings.Contains(file, "linux") {
							result += "linux"
						}
						if strings.Contains(file, "win") {
							result += "win"
						}
						if strings.Contains(file, "osx") {
							result += "darwin"
						}
						result += "-" + split[1]
						return
					}()),
					URL: DefaultURL + "v" + version + "/" + fmt.Sprintf("node-v%s-%s.%s", version, func() (result string) {

						if strings.Contains(file, "linux") {
							result += "linux"
						}
						if strings.Contains(file, "win") {
							result += "win"
						}
						if strings.Contains(file, "osx") {
							result += "darwin"
						}
						result += "-" + split[1]
						return
					}(), typeFile),
					Kind: func() string {
						if strings.Contains(file, "src") {
							return util.SourceKind
						}
						if strings.Contains(file, "msi") || strings.Contains(file, "pkg") {
							return util.InstallerKind
						}
						return util.ArchiveKind
					}(),
					OS: func() string {
						if strings.Contains(file, "linux") {
							return "linux"
						}
						if strings.Contains(file, "win") {
							return "windows"
						}
						if strings.Contains(file, "osx") {
							return "darwin"
						}
						return "unknown"
					}(),
					Arch:      split[1],
					Size:      "",
					Checksum:  "SHA256",
					Algorithm: "SHA256",
				})
			}
			meta[version] = nodeVersion
		}

		if element.Npm != "" {
			npm[version] = element.Npm
		}

		if isLTS(element) {
			lts = append(lts, version)
		} else if isCurrent(element) {
			current = append(current, version)
		} else if isStable(element) {
			stable = append(stable, version)
		} else if IsUnstable(element) {
			unstable = append(unstable, version)
		}
	}

	return all, lts, current, stable, unstable, npm, nil
}
