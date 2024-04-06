package commands_node

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/urfave/cli"
	"testing"
)

func TestCommandListInstalled(t *testing.T) {
	Convey("测试 标记", t, func() {

		CommandListInstalled(&cli.Context{})
	})
}

func TestCommandListRemote(t *testing.T) {
	Convey("测试线上版本拉取", t, func() {
		CommandListRemote(&cli.Context{})
	})
}

func TestCommandInstall(t *testing.T) {
	Convey("测试线上版本拉取", t, func() {

		err := commandInstall("21.7.2")
		if err != nil {
			t.Log(err)
		}
	})
}

func TestCommandListInstalled1(t *testing.T) {
	Convey("本地版本列表", t, func() {

		CommandListInstalled(&cli.Context{})

	})
}
