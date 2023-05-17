package commands_go

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
