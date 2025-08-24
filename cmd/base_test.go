package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBaseCommands(t *testing.T) {
	Convey("测试基础命令配置", t, func() {
		Convey("测试arch命令", func() {
			archCmd := baseCommands[0]
			So(archCmd.Name, ShouldEqual, "arch")
			So(archCmd.Usage, ShouldEqual, "systemc arch")
			So(archCmd.UsageText, ShouldEqual, "envm arch")
			So(archCmd.Action, ShouldNotBeNil)
		})

		Convey("测试go命令", func() {
			goCmd := baseCommands[1]
			So(goCmd.Name, ShouldEqual, "go")
			So(goCmd.Usage, ShouldEqual, "env go")
			So(goCmd.UsageText, ShouldEqual, "envm go")
			So(goCmd.Before, ShouldNotBeNil)
			So(goCmd.Subcommands, ShouldNotBeEmpty)
			So(len(goCmd.Subcommands), ShouldEqual, 5)
		})

		Convey("测试java命令", func() {
			javaCmd := baseCommands[2]
			So(javaCmd.Name, ShouldEqual, "java")
			So(javaCmd.Usage, ShouldEqual, "env java")
			So(javaCmd.UsageText, ShouldEqual, "envm java")
			So(javaCmd.Before, ShouldNotBeNil)
			So(javaCmd.Subcommands, ShouldNotBeEmpty)
			So(len(javaCmd.Subcommands), ShouldEqual, 3)
		})
	})
}

func TestGoSubcommands(t *testing.T) {
	Convey("测试Go子命令", t, func() {
		expectedCommands := []struct {
			name      string
			usage     string
			usageText string
		}{
			{"ls", "List installed versions", "envm ls"},
			{"lsr", "List remote versions available for install", "envm lsr [stable|archived]"},
			{"active", "Switch to specified version", "envm active <version>"},
			{"install", "Download and install a <version>", "envm install <version>"},
			{"uninstall", "Uninstall a version", "gvm uninstall <version>"},
		}

		So(len(goCommands), ShouldEqual, len(expectedCommands))

		for i, expected := range expectedCommands {
			cmd := goCommands[i]
			So(cmd.Name, ShouldEqual, expected.name)
			So(cmd.Usage, ShouldEqual, expected.usage)
			So(cmd.UsageText, ShouldEqual, expected.usageText)
			So(cmd.Action, ShouldNotBeNil)
		}
	})
}

func TestJavaSubcommands(t *testing.T) {
	Convey("测试Java子命令", t, func() {
		expectedCommands := []struct {
			name      string
			usage     string
			usageText string
		}{
			{"ls", "List installed versions", "envm java  ls"},
			{"active", "Switch to specified version", "envm java active <version>"},
			{"uninstall", "Uninstall a version", "envm java uninstall <version>"},
		}

		So(len(javaCommands), ShouldEqual, len(expectedCommands))

		for i, expected := range expectedCommands {
			cmd := javaCommands[i]
			So(cmd.Name, ShouldEqual, expected.name)
			So(cmd.Usage, ShouldEqual, expected.usage)
			So(cmd.UsageText, ShouldEqual, expected.usageText)
			So(cmd.Action, ShouldNotBeNil)
		}
	})
}