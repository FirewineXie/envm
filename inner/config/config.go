package config

import (
	"github.com/FirewineXie/govm/inner/arch"
	"os"
	"path"
	"path/filepath"
)

type GoVmConfig struct {
	Settings string `json:"settings"` // 配置文件  ${HOME}/.govm/settings
	Root     string `json:"root"`     // 路径 存储位置  ${HOME}/.govm/
	Symlink  string `json:"symlink"`  // 默认  ${HOME}/.govm/go
	Arch     string `json:"arch"`     // 系统arch
	Download string `json:"download"` // 默认为 ${HOME}/.govm/download   // 可以进行修改
}

var root = filepath.Clean(os.Getenv("GOVM_HOME"))
var symlink = filepath.Clean(os.Getenv("GOVM_SYMLINK"))

var env = GoVmConfig{
	Settings: path.Clean(path.Join(root, "settings")),
	Root:     root,
	Symlink:  symlink,
	Arch:     arch.Validate(""),
}

func Default() GoVmConfig {
	return env
}
