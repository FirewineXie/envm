package arch

import (
	"fmt"
	"github.com/urfave/cli"
)

func CommandArch(ctx *cli.Context) {
	fmt.Println(Validate())
}
