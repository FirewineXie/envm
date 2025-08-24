package cmd

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExecute(t *testing.T) {
	Convey("测试Execute函数", t, func() {
		Convey("测试version命令", func() {
			// 保存原始参数
			oldArgs := os.Args
			defer func() {
				os.Args = oldArgs
			}()

			// 模拟version命令
			os.Args = []string{"envm", "--version"}

			// Execute函数会调用os.Exit，所以我们不能直接测试它
			// 这里主要测试配置是否正确
			So(baseCommands, ShouldNotBeEmpty)
			So(len(baseCommands), ShouldEqual, 3) // arch, go, java
		})

		Convey("测试命令结构", func() {
			So(baseCommands[0].Name, ShouldEqual, "arch")
			So(baseCommands[1].Name, ShouldEqual, "go")
			So(baseCommands[2].Name, ShouldEqual, "java")
		})
	})
}