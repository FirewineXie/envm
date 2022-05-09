package arch

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"strings"
)

// Validate 通过系统变量，查看系统配置
func Validate(str string) string {
	if str == "" {
		str = os.Getenv("PROCESSOR_ARCHITECTURE")
	}

	return strings.ToLower(str)
}

func CommandArch(ctx *cli.Context) {
	fmt.Println(Validate(ctx.Args().First()))
}
