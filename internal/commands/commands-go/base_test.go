package commands_go

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
	Convey("测试列出已安装的Go版本", t, func() {
		// 创建测试CLI上下文
		app := cli.NewApp()
		ctx := cli.NewContext(app, nil, nil)

		// 执行命令
		err := CommandListInstalled(ctx)
		So(err, ShouldBeNil)
	})
}

func TestCommandListRemote(t *testing.T) {
	Convey("测试列出远程Go版本", t, func() {
		SkipConvey("跳过需要网络请求的测试", func() {
			app := cli.NewApp()
			ctx := cli.NewContext(app, nil, nil)

			err := CommandListRemote(ctx)
			So(err, ShouldBeNil)
		})
	})
}

func TestCommandUse(t *testing.T) {
	Convey("测试切换Go版本", t, func() {
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
	Convey("测试卸载Go版本", t, func() {
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

func TestCommandInstall(t *testing.T) {
	Convey("测试安装Go版本", t, func() {
		Convey("测试无参数情况", func() {
			app := cli.NewApp()
			ctx := cli.NewContext(app, nil, nil)

			err := CommandInstall(ctx)
			So(err, ShouldNotBeNil) // 应该返回错误
		})

		Convey("测试网络错误情况", func() {
			SkipConvey("跳过需要网络请求的测试", func() {})
		})
	})
}

func TestGoVersionUtilities(t *testing.T) {
	Convey("测试Go版本工具函数", t, func() {
		configLocal := config.Default().LinkSetting[config.GO]

		Convey("测试获取已安装版本", func() {
			if configLocal.Downloads != "" {
				versions := common.GetInstalled(configLocal.Downloads, "go")
				So(versions, ShouldNotBeNil)
				fmt.Printf("已安装的Go版本: %v\n", versions)
			} else {
				SkipConvey("ENVM_GO_SYMLINK 未配置", func() {})
			}
		})

		Convey("测试获取当前Go版本", func() {
			currentVersion := common.GetCurrentVersion("go")
			So(currentVersion, ShouldNotBeEmpty)
			fmt.Printf("当前Go版本: %s\n", currentVersion)
		})
	})
}

func TestGoConfigValidation(t *testing.T) {
	Convey("测试Go配置验证", t, func() {
		Convey("测试配置存在性", func() {
			configLocal := config.Default().LinkSetting[config.GO]
			
			if os.Getenv("ENVM_GO_SYMLINK") != "" {
				So(configLocal.Symlink, ShouldNotBeEmpty)
				So(configLocal.Downloads, ShouldNotBeEmpty)
			} else {
				SkipConvey("ENVM_GO_SYMLINK 环境变量未设置", func() {})
			}
		})

		Convey("测试Go环境验证", func() {
			err := config.VerifyEnvGo()
			if os.Getenv("ENVM_GO_SYMLINK") != "" {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "请先配置 ENVM_GO_SYMLINK")
			}
		})
	})
}

func TestListInstalledVersionsFormatting(t *testing.T) {
	Convey("测试已安装版本列表格式化", t, func() {
		configLocal := config.Default().LinkSetting[config.GO]
		
		if configLocal.Downloads != "" {
			currentVersion := common.GetCurrentVersion("go")
			installedVersions := common.GetInstalled(configLocal.Downloads, "go")
			
			So(currentVersion, ShouldNotBeEmpty)
			So(installedVersions, ShouldNotBeNil)
			
			// 测试格式化逻辑
			for _, version := range installedVersions {
				So(version, ShouldNotBeEmpty)
				// 版本应该符合语义版本格式
				So(version, ShouldNotStartWith, "go") // 已处理过的版本不应该有go前缀
			}
		} else {
			SkipConvey("Go环境未配置", func() {})
		}
	})
}