package commands_java

import (
	"fmt"
	"os"
	"testing"

	"github.com/FirewineXie/envm/internal/commands/common"
	"github.com/FirewineXie/envm/internal/config"
	"github.com/urfave/cli"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCommandListInstalled(t *testing.T) {
	Convey("测试列出已安装的Java版本", t, func() {
		// 创建测试CLI上下文
		app := cli.NewApp()
		ctx := cli.NewContext(app, nil, nil)

		// 执行命令
		err := CommandListInstalled(ctx)
		So(err, ShouldBeNil)
	})
}

func TestCommandUse(t *testing.T) {
	Convey("测试切换Java版本", t, func() {
		Convey("测试无参数情况", func() {
			app := cli.NewApp()
			ctx := cli.NewContext(app, nil, nil)

			// 测试无参数时的行为
			err := CommandUse(ctx)
			So(err, ShouldNotBeNil) // 应该返回错误或帮助信息
		})

		Convey("测试版本不存在", func() {
			SkipConvey("需要mock本地版本检查", func() {})
		})
	})
}

func TestCommandUninstall(t *testing.T) {
	Convey("测试卸载Java版本", t, func() {
		Convey("测试无参数情况", func() {
			app := cli.NewApp()
			ctx := cli.NewContext(app, nil, nil)

			err := CommandUninstall(ctx)
			// 应该有适当的错误处理
			So(err, ShouldNotBeNil)
		})

		Convey("测试卸载当前版本", func() {
			SkipConvey("需要mock当前版本检查", func() {})
		})
	})
}

func TestJavaVersionUtilities(t *testing.T) {
	Convey("测试Java版本工具函数", t, func() {
		configLocal := config.Default().LinkSetting[config.JAVA]

		Convey("测试获取已安装版本", func() {
			if configLocal.Downloads != "" {
				versions := common.GetInstalled(configLocal.Downloads, "jdk")
				So(versions, ShouldNotBeNil)
				fmt.Printf("已安装的Java版本: %v\n", versions)
			} else {
				SkipConvey("ENVM_JAVA_SYMLINK 未配置", func() {})
			}
		})

		Convey("测试获取当前Java版本", func() {
			currentVersion := common.GetCurrentVersion("java")
			So(currentVersion, ShouldNotBeEmpty)
			fmt.Printf("当前Java版本: %s\n", currentVersion)
		})
	})
}

func TestJavaConfigValidation(t *testing.T) {
	Convey("测试Java配置验证", t, func() {
		Convey("测试配置存在性", func() {
			configLocal := config.Default().LinkSetting[config.JAVA]
			
			if os.Getenv("ENVM_JAVA_SYMLINK") != "" {
				So(configLocal.Symlink, ShouldNotBeEmpty)
				So(configLocal.Downloads, ShouldNotBeEmpty)
			} else {
				SkipConvey("ENVM_JAVA_SYMLINK 环境变量未设置", func() {})
			}
		})

		Convey("测试Java环境验证", func() {
			err := config.VerifyEnvJava()
			if os.Getenv("ENVM_JAVA_SYMLINK") != "" {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "请先配置 ENVM_JAVA_SYMLINK")
			}
		})
	})
}

func TestJavaListInstalledVersionsFormatting(t *testing.T) {
	Convey("测试Java已安装版本列表格式化", t, func() {
		configLocal := config.Default().LinkSetting[config.JAVA]
		
		if configLocal.Downloads != "" {
			currentVersion := common.GetCurrentVersion("java")
			installedVersions := common.GetInstalled(configLocal.Downloads, "jdk")
			
			So(currentVersion, ShouldNotBeEmpty)
			So(installedVersions, ShouldNotBeNil)
			
			// 测试格式化逻辑
			for _, version := range installedVersions {
				So(version, ShouldNotBeEmpty)
				// 版本应该符合语义版本格式
				So(version, ShouldNotStartWith, "jdk") // 已处理过的版本不应该有jdk前缀
			}
		} else {
			SkipConvey("Java环境未配置", func() {})
		}
	})
}