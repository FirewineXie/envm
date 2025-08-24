package config

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVerifyEnv(t *testing.T) {
	Convey("测试环境变量验证", t, func() {
		Convey("测试ENVM_HOME为空时", func() {
			// 保存原始环境变量
			originalRoot := root
			originalArch := env.Arch
			defer func() {
				root = originalRoot
				env.Root = originalRoot
				env.Arch = originalArch
			}()

			// 模拟空root
			root = "."
			env.Root = "."

			err := VerifyEnv()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "root 路径不能为空")
		})

		Convey("测试arch为空时", func() {
			// 保存原始环境变量
			originalRoot := root
			originalArch := env.Arch
			defer func() {
				root = originalRoot
				env.Root = originalRoot
				env.Arch = originalArch
			}()

			// 模拟有效root但空arch
			root = "/tmp/envm"
			env.Root = "/tmp/envm"
			env.Arch = ""

			err := VerifyEnv()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "arch 暂时不支持")
		})

		Convey("测试正常情况", func() {
			// 保存原始环境变量
			originalRoot := root
			originalArch := env.Arch
			defer func() {
				root = originalRoot
				env.Root = originalRoot
				env.Arch = originalArch
			}()

			// 模拟正常配置
			root = "/tmp/envm"
			env.Root = "/tmp/envm"
			env.Arch = "amd64"

			err := VerifyEnv()
			So(err, ShouldBeNil)
		})
	})
}

func TestVerifyEnvGo(t *testing.T) {
	Convey("测试Go环境变量验证", t, func() {
		Convey("测试ENVM_GO_SYMLINK未配置", func() {
			// 保存原始配置
			originalLinkSetting := env.LinkSetting
			defer func() {
				env.LinkSetting = originalLinkSetting
			}()

			// 模拟空配置
			env.LinkSetting = map[string]SubConfig{
				GO: {Symlink: ""},
			}

			err := VerifyEnvGo()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "请先配置 ENVM_GO_SYMLINK")
		})

		Convey("测试ENVM_GO_SYMLINK已配置", func() {
			// 保存原始配置
			originalLinkSetting := env.LinkSetting
			defer func() {
				env.LinkSetting = originalLinkSetting
			}()

			// 模拟正常配置
			env.LinkSetting = map[string]SubConfig{
				GO: {Symlink: "/tmp/envm/go"},
			}

			err := VerifyEnvGo()
			So(err, ShouldBeNil)
		})
	})
}

func TestVerifyEnvJava(t *testing.T) {
	Convey("测试Java环境变量验证", t, func() {
		Convey("测试ENVM_JAVA_SYMLINK未配置", func() {
			// 保存原始配置
			originalLinkSetting := env.LinkSetting
			defer func() {
				env.LinkSetting = originalLinkSetting
			}()

			// 模拟空配置
			env.LinkSetting = map[string]SubConfig{
				JAVA: {Symlink: ""},
			}

			err := VerifyEnvJava()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "请先配置 ENVM_JAVA_SYMLINK")
		})

		Convey("测试ENVM_JAVA_SYMLINK已配置", func() {
			// 保存原始配置
			originalLinkSetting := env.LinkSetting
			defer func() {
				env.LinkSetting = originalLinkSetting
			}()

			// 模拟正常配置
			env.LinkSetting = map[string]SubConfig{
				JAVA: {Symlink: "/tmp/envm/java"},
			}

			err := VerifyEnvJava()
			So(err, ShouldBeNil)
		})
	})
}

func TestDefault(t *testing.T) {
	Convey("测试默认配置", t, func() {
		config := Default()

		So(config.Root, ShouldNotBeEmpty)
		So(config.Arch, ShouldNotBeEmpty)
		So(config.Downloads, ShouldNotBeEmpty)
		So(config.LinkSetting, ShouldNotBeNil)
	})
}

func TestEnvmConfig(t *testing.T) {
	Convey("测试EnvmConfig结构", t, func() {
		Convey("测试配置常量", func() {
			So(GO, ShouldEqual, "go")
			So(JAVA, ShouldEqual, "java")
		})

		Convey("测试SubConfig结构", func() {
			subConfig := SubConfig{
				Symlink:   "/path/to/symlink",
				Downloads: "/path/to/downloads",
			}
			So(subConfig.Symlink, ShouldEqual, "/path/to/symlink")
			So(subConfig.Downloads, ShouldEqual, "/path/to/downloads")
		})
	})
}

func TestConfigInitialization(t *testing.T) {
	Convey("测试配置初始化", t, func() {
		// 测试当前环境下的初始化
		config := Default()

		// 验证基本配置
		So(config.Root, ShouldEqual, filepath.Clean(os.Getenv("ENVM_HOME")))
		So(config.Arch, ShouldNotBeEmpty)
		So(config.Downloads, ShouldEqual, filepath.Join(config.Root, "downloads"))

		// 如果环境变量存在，验证链接配置
		if os.Getenv("ENVM_GO_SYMLINK") != "" {
			goConfig, exists := config.LinkSetting[GO]
			So(exists, ShouldBeTrue)
			So(goConfig.Symlink, ShouldEqual, filepath.Clean(os.Getenv("ENVM_GO_SYMLINK")))
		}

		if os.Getenv("ENVM_JAVA_SYMLINK") != "" {
			javaConfig, exists := config.LinkSetting[JAVA]
			So(exists, ShouldBeTrue)
			So(javaConfig.Symlink, ShouldEqual, filepath.Clean(os.Getenv("ENVM_JAVA_SYMLINK")))
		}
	})
}