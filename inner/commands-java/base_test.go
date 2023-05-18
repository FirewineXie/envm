package commands_java

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
