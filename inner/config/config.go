package config

import (
	"github.com/FirewineXie/govm/inner/arch"
	"os"
	"path/filepath"
)

type GoVmConfig struct {
	SettingPath string `json:"setting_path"` // 配置文件名称
	Root        string `json:"root"`         // 执行命令位置
	Symlink     string `json:"symlink"`      // 链接位置
	Arch        string `json:"arch"`         // 系统arch
	Downloads   string `json:"downloads"`    // 下载目录
}

var root = filepath.Clean(os.Getenv("GOVM_HOME"))
var symlink = filepath.Clean(os.Getenv("GOVM_SYMLINK"))

var env = GoVmConfig{
	SettingPath: filepath.Join(root, "settings"),
	Root:        root,
	Symlink:     symlink,
	Arch:        arch.Validate(),
	Downloads:   filepath.Join(root, "downloads"),
}

func Default() GoVmConfig {
	ReadSettings()
	return env
}

func VerifyEnv() bool {
	if root == "settings" {
		return false
	}
	if symlink == "" {
		return false
	}
	if env.Arch == "" {
		return false
	}
	return true
}
