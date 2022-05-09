package config

import (
	"github.com/FirewineXie/govm/inner/arch"
	"os"
	"path/filepath"
)

type GoVmConfig struct {
	Settings string `json:"settings"` // 配置文件
	Proxy    string `json:"proxy"`    // 代理
	Root     string `json:"root"`     // 路径
	Symlink  string `json:"symlink"`  // 系统链接
	Arch     string `json:"arch"`     // 系统arch
}

var home = filepath.Clean(os.Getenv("GOVM_HOME") + "/settings")

var env = GoVmConfig{
	Settings: home,
	Root:     "",
	Symlink:  "",
	Arch:     arch.Validate(""),
	Proxy:    "none",
}

func Default() GoVmConfig {
	return env
}

func WithSymlink(symlink string) GoVmConfig {
	env.Symlink = symlink
	return env
}
