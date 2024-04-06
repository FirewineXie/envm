package config

import (
	"errors"
	"github.com/FirewineXie/envm/internal/arch"
	"github.com/FirewineXie/envm/util"

	"os"
	"path/filepath"
)

type EnvmConfig struct {
	Root        string `json:"root"`      // root 目录
	Arch        string `json:"arch"`      // 系统arch
	Downloads   string `json:"downloads"` // 下载目录
	LinkSetting map[string]SubConfig
	Settings    Settings `json:"settings"`
}

type Settings struct {
}

type SubConfig struct {
	Symlink   string `json:"symlink"`   // 链接位置
	Downloads string `json:"downloads"` // 相对位置
}

var root = filepath.Clean(os.Getenv("ENVM_HOME"))
var goSymlink = filepath.Clean(os.Getenv("ENVM_GO_SYMLINK"))
var javaSymlink = filepath.Clean(os.Getenv("ENVM_JAVA_SYMLINK"))
var nodeSymlink = filepath.Clean(os.Getenv("ENVM_NODE_SYMLINK"))

var env = EnvmConfig{
	Root:        root,
	Arch:        arch.Validate(),
	LinkSetting: map[string]SubConfig{},
}

func init() {
	env.Downloads = filepath.Join(root, "downloads")
	exists, _ := util.PathExists(env.Downloads)
	if !exists {
		_ = os.Mkdir(env.Downloads, os.ModePerm)
	}

	if goSymlink != "." {
		env.LinkSetting[GO] = SubConfig{
			goSymlink,
			filepath.Join(env.Downloads, "go"),
		}
		pathExists, _ := util.PathExists(env.LinkSetting[GO].Downloads)
		if !pathExists {
			_ = os.Mkdir(env.LinkSetting[GO].Downloads, os.ModePerm)
		}

	}

	if javaSymlink != "." {
		env.LinkSetting[JAVA] = SubConfig{
			javaSymlink,
			filepath.Join(env.Downloads, "java"),
		}
		pathExists, _ := util.PathExists(env.LinkSetting[JAVA].Downloads)
		if !pathExists {
			_ = os.Mkdir(env.LinkSetting[JAVA].Downloads, os.ModePerm)
		}
	}
	if javaSymlink != "." {
		env.LinkSetting[NODE] = SubConfig{
			nodeSymlink,
			filepath.Join(env.Downloads, "node"),
		}
		pathExists, _ := util.PathExists(env.LinkSetting[NODE].Downloads)
		if !pathExists {
			_ = os.Mkdir(env.LinkSetting[NODE].Downloads, os.ModePerm)
		}
	}
}

const (
	GO   = "go"
	JAVA = "java"
	NODE = "node"
)

func Default() EnvmConfig {
	return env
}

func VerifyEnv() error {
	if root == "." {
		return errors.New("root 路径不能为空，请配置  ENVM_HOME 为当前执行程序路径")
	}
	if env.Arch == "" {
		return errors.New("arch 暂时不支持")
	}

	return nil
}

func VerifyEnvGo() error {
	symlink := env.LinkSetting[GO].Symlink

	if symlink == "" {
		return errors.New("请先配置 ENVM_GO_SYMLINK")
	}
	return nil
}

func VerifyEnvJava() error {
	symlink := env.LinkSetting[JAVA].Symlink

	if symlink == "" {
		return errors.New("请先配置 ENVM_JAVA_SYMLINK")
	}
	return nil
}
func VerifyEnvNode() error {
	symlink := env.LinkSetting[NODE].Symlink

	if symlink == "" {
		return errors.New("请先配置 ENVM_NODE_SYMLINK")
	}
	return nil
}
