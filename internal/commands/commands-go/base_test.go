package commands_go

import (
	"fmt"
	"github.com/FirewineXie/envm/internal/commands/common"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/urfave/cli"
	"regexp"
	"testing"
)

func TestCommandListInstalled(t *testing.T) {
	Convey("测试 标记", t, func() {

		CommandListInstalled(&cli.Context{})
	})
}

func TestCommandListRemote(t *testing.T) {
	in := common.GetCurrentVersion("go")

	v := common.GetInstalled(configLocal.Downloads, "go")

	for i := 0; i < len(v); i++ {
		version := v[i]

		str := ""
		goVersion := fmt.Sprintf("go%v", version)
		if in == goVersion {
			str = str + "  * "
		} else {
			str = str + "    "
		}
		str = str + regexp.MustCompile("go").ReplaceAllString(version, "")
		if in == goVersion {
			str = str + " (Currently using " + in + " executable)"
		}
		fmt.Printf(str + "\n")

	}
	if len(v) == 0 {
		fmt.Println("No installations recognized.")
	}
}
