//go:build windows

package arch

import (
	"os"
	"strings"
)

// Validate 通过系统变量，查看系统配置
func Validate() string {

	str := os.Getenv("PROCESSOR_ARCHITECTURE")

	return strings.ToLower(str)
}

func GetArch() (string, error) {
	str := os.Getenv("PROCESSOR_ARCHITECTURE")

	return strings.ToLower(str), nil
}
