package util

import (
	"runtime"
)

// CreateSymlink 创建符号链接，兼容Windows和Unix系统
// Windows版本会自动处理UAC权限提升
func CreateSymlink(target, linkPath string) error {
	if runtime.GOOS == "windows" {
		// Windows上使用UAC权限提升功能
		return CreateSymlinkWithElevation(target, linkPath)
	}
	
	// Unix系统使用跨平台函数
	return CreateSymlinkWithElevation(target, linkPath)
}

// WindowsSymlinkError Windows符号链接错误（保留向后兼容）
type WindowsSymlinkError struct {
	Target   string
	LinkPath string
	Err      error
}

func (e *WindowsSymlinkError) Error() string {
	return "Failed to create symlink on Windows. UAC elevation may be required.\n" +
		"Original error: " + e.Err.Error()
}