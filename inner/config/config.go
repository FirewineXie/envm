package config

import (
	"github.com/FirewineXie/govm/inner/arch"
	"github.com/FirewineXie/govm/util"
	"os"
	"path/filepath"
)

type GoVmConfig struct {
	Root      string `json:"root"`      // 执行命令位置
	Symlink   string `json:"symlink"`   // 链接位置
	Arch      string `json:"arch"`      // 系统arch
	Downloads string `json:"downloads"` // 下载目录
}

var root = filepath.Clean(os.Getenv("GOVM_HOME"))
var symlink = filepath.Clean(os.Getenv("GOVM_SYMLINK"))
var downloads = filepath.Clean(os.Getenv("GOVM_DOWNLOAD"))

var env = GoVmConfig{
	Root:      root,
	Symlink:   symlink,
	Arch:      arch.Validate(),
	Downloads: downloads,
}

func Default() GoVmConfig {
	if env.Downloads == "" {
		env.Downloads = filepath.Join(root, "downloads")
	}
	exists, _ := util.PathExists(env.Downloads)
	if !exists {
		_ = os.Mkdir(env.Downloads, os.ModePerm)
	}

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
	if env.Downloads == "" {
		return false
	}
	return true
}
