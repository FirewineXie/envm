//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package arch

import (
	"os/exec"
	"strings"
)

// Validate 通过系统变量，查看系统配置
func Validate() string {

	arch, _ := GetArch()

	return strings.ToLower(strings.Replace(arch, "\n", "", -1))

}

func GetArch() (string, error) {
	command := exec.Command("/bin/bash", "-c", `uname -m`)
	output, err := command.CombinedOutput()
	return string(output), err
}
