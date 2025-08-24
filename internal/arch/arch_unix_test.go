package arch

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/urfave/cli"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidate(t *testing.T) {
	Convey("测试架构验证", t, func() {
		arch := Validate()
		So(arch, ShouldNotBeEmpty)
		
		// 验证返回的架构是否为预期的架构之一
		expectedArchs := []string{"amd64", "386", "arm", "arm64"}
		So(expectedArchs, ShouldContain, arch)
		
		// 验证与runtime.GOARCH一致
		So(arch, ShouldEqual, runtime.GOARCH)
		
		fmt.Println("当前系统架构:", arch)
	})
}

func TestCommandArch(t *testing.T) {
	Convey("测试arch命令", t, func() {
		// 创建一个CLI应用用于测试
		app := cli.NewApp()
		app.Commands = []cli.Command{
			{
				Name:   "arch",
				Action: CommandArch,
			},
		}

		// 捕获输出
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// 创建测试context
		ctx := cli.NewContext(app, nil, nil)

		// 执行命令
		CommandArch(ctx)

		// 恢复输出并读取结果
		w.Close()
		os.Stdout = oldStdout
		
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// 验证输出
		So(output, ShouldNotBeEmpty)
		So(output, ShouldContainSubstring, runtime.GOARCH)
	})
}

func TestArchitectureSupport(t *testing.T) {
	Convey("测试不同架构支持", t, func() {
		arch := Validate()
		
		Convey("测试AMD64架构", func() {
			if runtime.GOARCH == "amd64" {
				So(arch, ShouldEqual, "amd64")
			}
		})
		
		Convey("测试ARM64架构", func() {
			if runtime.GOARCH == "arm64" {
				So(arch, ShouldEqual, "arm64")
			}
		})
		
		Convey("测试386架构", func() {
			if runtime.GOARCH == "386" {
				So(arch, ShouldEqual, "386")
			}
		})
		
		Convey("测试ARM架构", func() {
			if runtime.GOARCH == "arm" {
				So(arch, ShouldEqual, "arm")
			}
		})
	})
}
