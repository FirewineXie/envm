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
