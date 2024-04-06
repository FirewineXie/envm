package web_node

import "github.com/blang/semver/v4"

/*
 * @Author: Firewine
 * @File: node
 * @Version: 1.0.0
 * @Date: 2024-04-05 17:19
 * @Description:
 */

// isLTS Identifies a version as "LTS"
func isLTS(element FileData) bool {
	switch dataType := element.Lts.(type) {
	case bool:
		return dataType
	case string:
		return true
	}
	return false
}

// isCurrent Identifies a version as "current"
func isCurrent(element FileData) bool {
	if isLTS(element) {
		return false
	}

	version, _ := semver.Make(element.Version[1:])
	benchmark, _ := semver.Make("1.0.0")

	if version.LT(benchmark) {
		return false
	}

	return true
}

// isStable Identifies a stable old version.
func isStable(element FileData) bool {
	if isCurrent(element) {
		return false
	}

	version, _ := semver.Make(element.Version[1:])

	if version.Major != 0 {
		return false
	}

	return version.Minor%2 == 0
}

// IsUnstable Identifies an unstable old version.
func IsUnstable(element FileData) bool {
	if isStable(element) {
		return false
	}

	version, _ := semver.Make(element.Version[1:])

	if version.Major != 0 {
		return false
	}

	return version.Minor%2 != 0
}

func exchangeArch(arch string) string {
	if arch == "amd64" {
		return "x64"
	}
	return arch
}

type FileData struct {
	Version  string   `json:"version"`
	Date     string   `json:"date"`
	Files    []string `json:"files"`
	Npm      string   `json:"npm,omitempty"`
	V8       string   `json:"v8"`
	Uv       string   `json:"uv,omitempty"`
	Zlib     string   `json:"zlib,omitempty"`
	Openssl  string   `json:"openssl,omitempty"`
	Modules  string   `json:"modules,omitempty"`
	Lts      any      `json:"lts"`
	Security bool     `json:"security"`
}
