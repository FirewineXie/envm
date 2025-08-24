package util

import (
	"errors"
)

// ErrVersionNotFound 版本不存在
var ErrVersionNotFound = errors.New("version not found")

// FindVersion 返回指定名称的版本
func FindVersion(all []*Version, name string) (*Version, error) {
	for i := range all {
		if all[i].Name == name {
			return all[i], nil
		}
	}
	return nil, ErrVersionNotFound
}

// FindVersionV2 返回指定名称的版本
func FindVersionV2(all map[string]*Version, name string) (*Version, error) {
	if v, ok := all[name]; ok {
		return v, nil
	}
	return nil, ErrVersionNotFound
}

// Version 版本
type Version struct {
	Name     string // 版本名，如'1.12.4'
	Packages []*Package
}

// ErrPackageNotFound 版本包不存在
var ErrPackageNotFound = errors.New("installation package not found")

type FindPackageInterface interface {
	FindPackage(king, goos, arch string) (*Package, error)
}

// FindPackage 查找指定类型、操作系统和架构的包
func (v *Version) FindPackage(kind, goos, arch string) (*Package, error) {
	for i := range v.Packages {
		pkg := v.Packages[i]
		if pkg.Kind == kind {
			// 匹配操作系统
			osMatch := false
			switch goos {
			case "darwin":
				osMatch = pkg.OS == "macOS" || pkg.OS == "Darwin"
			case "windows":
				osMatch = pkg.OS == "Windows"
			case "linux":
				osMatch = pkg.OS == "Linux"
			default:
				osMatch = pkg.OS == goos
			}
			
			// 匹配架构
			archMatch := false
			switch arch {
			case "amd64":
				archMatch = pkg.Arch == "x86-64" || pkg.Arch == "amd64"
			case "386":
				archMatch = pkg.Arch == "x86" || pkg.Arch == "386"
			case "arm64":
				archMatch = pkg.Arch == "ARM64" || pkg.Arch == "arm64"
			default:
				archMatch = pkg.Arch == arch
			}
			
			if osMatch && archMatch {
				return pkg, nil
			}
		}
	}
	return nil, ErrPackageNotFound
}
