package config

import (
	"errors"
	"github.com/FirewineXie/govm/inner/arch"
	"github.com/FirewineXie/govm/util"
	"os"
	"path/filepath"
)

type EnvmConfig struct {
	Root      string               `json:"root"`      // root 目录
	Arch      string               `json:"arch"`      // 系统arch
	Downloads string               `json:"downloads"` // 下载目录
	Settings  map[string]SubConfig `json:"settings"`
}

type SubConfig struct {
	Symlink   string `json:"symlink"`   // 链接位置
	Downloads string `json:"downloads"` // 相对位置
}

var root = filepath.Clean(os.Getenv("ENVM_HOME"))
var goSymlink = filepath.Clean(os.Getenv("ENVM_GO_SYMLINK"))
var javaSymlink = filepath.Clean(os.Getenv("ENVM_JAVA_SYMLINK"))

var env = EnvmConfig{
	Root:     root,
	Arch:     arch.Validate(),
	Settings: map[string]SubConfig{},
}

func init() {
	env.Downloads = filepath.Join(root, "downloads")
	exists, _ := util.PathExists(env.Downloads)
	if !exists {
		_ = os.Mkdir(env.Downloads, os.ModePerm)
	}

	if goSymlink != "." {
		env.Settings[GO] = SubConfig{
			goSymlink,
			filepath.Join(env.Downloads, "go"),
		}
		pathExists, _ := util.PathExists(env.Settings[GO].Downloads)
		if !pathExists {
			_ = os.Mkdir(env.Settings[GO].Downloads, os.ModePerm)
		}

	}

	if javaSymlink != "." {
		env.Settings[JAVA] = SubConfig{
			javaSymlink,
			filepath.Join(env.Downloads, "java"),
		}
		pathExists, _ := util.PathExists(env.Settings[JAVA].Downloads)
		if !pathExists {
			_ = os.Mkdir(env.Settings[JAVA].Downloads, os.ModePerm)
		}
	}

}

const (
	GO   = "go"
	JAVA = "java"
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
	symlink := env.Settings[GO].Symlink

	if symlink == "" {
		return errors.New("请先配置 ENVM_GO_SYMLINK")
	}
	return nil
}

func VerifyEnvJava() error {
	symlink := env.Settings[JAVA].Symlink

	if symlink == "" {
		return errors.New("请先配置 ENVM_JAVA_SYMLINK")
	}
	return nil
}
