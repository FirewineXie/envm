//go:build windows

package util

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	advapi32                = syscall.NewLazyDLL("advapi32.dll")
	shell32                 = syscall.NewLazyDLL("shell32.dll")
	procGetCurrentProcess   = kernel32.NewProc("GetCurrentProcess")
	procOpenProcessToken    = advapi32.NewProc("OpenProcessToken")
	procGetTokenInformation = advapi32.NewProc("GetTokenInformation")
	procShellExecuteW       = shell32.NewProc("ShellExecuteW")
)

const (
	TOKEN_QUERY         = 0x0008
	TokenElevation      = 20
	TokenElevationType  = 18
	TokenElevationTypeLimited = 3
)

// IsRunningAsAdmin 检查当前进程是否以管理员权限运行
func IsRunningAsAdmin() bool {
	var token syscall.Handle
	var elevation uint32
	var elevationSize uint32

	// 获取当前进程令牌
	currentProcess, _, _ := procGetCurrentProcess.Call()
	ret, _, _ := procOpenProcessToken.Call(
		currentProcess,
		TOKEN_QUERY,
		uintptr(unsafe.Pointer(&token)),
	)
	if ret == 0 {
		return false
	}
	defer syscall.CloseHandle(token)

	// 检查令牌提升状态
	ret, _, _ = procGetTokenInformation.Call(
		uintptr(token),
		TokenElevation,
		uintptr(unsafe.Pointer(&elevation)),
		uintptr(unsafe.Sizeof(elevation)),
		uintptr(unsafe.Pointer(&elevationSize)),
	)
	if ret == 0 {
		return false
	}

	return elevation != 0
}

// CanElevatePermissions 检查当前用户是否可以提升权限
func CanElevatePermissions() bool {
	var token syscall.Handle
	var elevationType uint32
	var elevationTypeSize uint32

	// 获取当前进程令牌
	currentProcess, _, _ := procGetCurrentProcess.Call()
	ret, _, _ := procOpenProcessToken.Call(
		currentProcess,
		TOKEN_QUERY,
		uintptr(unsafe.Pointer(&token)),
	)
	if ret == 0 {
		return false
	}
	defer syscall.CloseHandle(token)

	// 检查令牌提升类型
	ret, _, _ = procGetTokenInformation.Call(
		uintptr(token),
		TokenElevationType,
		uintptr(unsafe.Pointer(&elevationType)),
		uintptr(unsafe.Sizeof(elevationType)),
		uintptr(unsafe.Pointer(&elevationTypeSize)),
	)
	if ret == 0 {
		return false
	}

	// TokenElevationTypeLimited 表示用户在管理员组中但当前权限受限
	return elevationType == TokenElevationTypeLimited
}

// RequestAdminPrivileges 请求管理员权限，显示UAC弹窗
func RequestAdminPrivileges(targetPath, linkPath string) error {
	if IsRunningAsAdmin() {
		// 已经是管理员权限，直接创建符号链接
		return os.Symlink(targetPath, linkPath)
	}

	if !CanElevatePermissions() {
		return fmt.Errorf("当前用户不在管理员组中，无法提升权限。请联系系统管理员或启用开发者模式")
	}

	// 获取当前可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取当前程序路径: %v", err)
	}

	// 构造参数
	args := fmt.Sprintf("--admin-symlink %s %s", targetPath, linkPath)
	
	// 使用ShellExecute调用UAC提升
	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exePath)
	params, _ := syscall.UTF16PtrFromString(args)
	
	ret, _, _ := procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(params)),
		0,
		1, // SW_SHOWNORMAL
	)

	if ret <= 32 {
		return fmt.Errorf("UAC权限请求失败，错误代码: %d", ret)
	}

	return nil
}

// CreateSymlinkWithElevation 创建符号链接，必要时请求权限提升
func CreateSymlinkWithElevation(targetPath, linkPath string) error {
	// 删除已存在的链接
	_ = os.Remove(linkPath)

	// 首先尝试直接创建符号链接
	err := os.Symlink(targetPath, linkPath)
	if err == nil {
		return nil
	}

	// 如果失败，检查是否是权限问题
	if !strings.Contains(err.Error(), "privilege") && 
	   !strings.Contains(err.Error(), "A required privilege is not held by the client") {
		return err // 不是权限问题，返回原错误
	}

	// 尝试权限提升
	fmt.Println("需要管理员权限来创建符号链接...")
	fmt.Println("将弹出UAC权限请求窗口，请点击'是'以授权。")
	
	return RequestAdminPrivileges(targetPath, linkPath)
}

// HandleAdminSymlinkCommand 处理管理员权限下的符号链接创建
func HandleAdminSymlinkCommand(args []string) error {
	if len(args) != 3 || args[0] != "--admin-symlink" {
		return fmt.Errorf("无效的管理员符号链接命令参数")
	}

	targetPath := args[1]
	linkPath := args[2]

	if !IsRunningAsAdmin() {
		return fmt.Errorf("此命令需要管理员权限")
	}

	// 删除已存在的链接
	_ = os.Remove(linkPath)

	// 创建符号链接
	err := os.Symlink(targetPath, linkPath)
	if err != nil {
		return fmt.Errorf("创建符号链接失败: %v", err)
	}

	fmt.Printf("符号链接创建成功: %s -> %s\n", linkPath, targetPath)
	return nil
}

// IsAdminSymlinkCommand 检查是否是管理员符号链接命令
func IsAdminSymlinkCommand(args []string) bool {
	return len(args) >= 1 && args[0] == "--admin-symlink"
}