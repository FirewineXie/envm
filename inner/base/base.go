package base

import (
	"fmt"
	"github.com/FirewineXie/govm/inner/config"
	"github.com/FirewineXie/govm/inner/web"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// CommandUse 激活使用go版本
func CommandUse(ctx *cli.Context) error {
	v, err := getVersion(ctx, true)
	if err != nil {

		return err
	}
	// active use
	_ = os.Remove(config.Default().Root)

	if err := os.Symlink(v, config.Default().Root); err != nil {
		return cli.NewExitError(fmt.Sprintf("%s", err.Error()), 1)
	}
	if output, err := exec.Command(filepath.Join(config.Default().Root, "bin", "go"), "version").Output(); err == nil {
		fmt.Printf(string(output))
	}
}

// CommandListRemote 获取远程的可下载的版本
func CommandListRemote(ctx *cli.Context) {
	collector, err := web.NewCollector("")
	if err != nil {
		log.Fatal(err)
		return
	}

	versions, err := collector.AllVersions()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(versions)
}

// CommandListInstalled 展示已经安装的go 版本
func CommandListInstalled(ctx *cli.Context) {
	in := getCurrentVersion()

	v := getInstalled(config.Default().Root)

	for i := 0; i < len(v); i++ {
		version := v[i]

		str := ""
		goVersion := fmt.Sprintf("go%v", version)
		if in == goVersion {
			str = str + "  * "
		} else {
			str = str + "    "
		}
		if in == goVersion {
			str = str + " (Currently using " + in + " executable)"
		}
		fmt.Printf(str + "\n")

	}
	if len(v) == 0 {
		fmt.Println("No installations recognized.")
	}
}
