//go:build !windows

package util

import "os"

// IsRunningAsAdmin Unix系统下的权限检查（通常不需要特殊处理）
func IsRunningAsAdmin() bool {
	return os.Geteuid() == 0
}

// CanElevatePermissions Unix系统下检查是否可以提升权限
func CanElevatePermissions() bool {
	return true // Unix系统通常可以使用sudo
}

// CreateSymlinkWithElevation Unix系统下创建符号链接
func CreateSymlinkWithElevation(targetPath, linkPath string) error {
	// 删除已存在的链接
	_ = os.Remove(linkPath)
	
	// Unix系统直接创建符号链接
	return os.Symlink(targetPath, linkPath)
}

// HandleAdminSymlinkCommand Unix系统下不需要特殊处理
func HandleAdminSymlinkCommand(args []string) error {
	return nil
}

// IsAdminSymlinkCommand Unix系统下总是返回false
func IsAdminSymlinkCommand(args []string) bool {
	return false
}