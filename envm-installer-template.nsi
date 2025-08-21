; NSIS 安装脚本模板 for envm
; 这个文件会被 GitHub Actions 动态替换变量

; 定义安装程序信息（会被 GitHub Actions 替换）
!define PRODUCT_NAME "envm"
!define PRODUCT_VERSION "{{VERSION}}"
!define PRODUCT_PUBLISHER "FirewineXie"
!define PRODUCT_WEB_SITE "https://github.com/FirewineXie/envm"
!define PRODUCT_DIR_REGKEY "Software\Microsoft\Windows\CurrentVersion\App Paths\envm.exe"
!define PRODUCT_UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"
!define PRODUCT_UNINST_ROOT_KEY "HKCU"
!define BINARY_FILE "{{BINARY_FILE}}"

; 设置安装程序属性
Name "${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "envm-installer-${PRODUCT_VERSION}.exe"
InstallDir "$LOCALAPPDATA\${PRODUCT_NAME}"
InstallDirRegKey HKCU "${PRODUCT_DIR_REGKEY}" ""
ShowInstDetails show
ShowUnInstDetails show
RequestExecutionLevel user

; 现代界面设置
!include "MUI2.nsh"
!include "LogicLib.nsh"
!include "FileFunc.nsh"
!define MUI_ABORTWARNING
!define MUI_ICON "${NSISDIR}\Contrib\Graphics\Icons\modern-install.ico"
!define MUI_UNICON "${NSISDIR}\Contrib\Graphics\Icons\modern-uninstall.ico"

; 自定义安装目录页面设置
!define MUI_DIRECTORYPAGE_TEXT_TOP "选择 ${PRODUCT_NAME} 的安装位置。"
!define MUI_DIRECTORYPAGE_TEXT_DESTINATION "安装文件夹"
!define MUI_DIRECTORYPAGE_VERIFYONLEAVE

; 安装页面
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!define MUI_FINISHPAGE_RUN "$INSTDIR\envm.exe"
!define MUI_FINISHPAGE_RUN_PARAMETERS "--version"
!define MUI_FINISHPAGE_RUN_TEXT "运行 envm --version"
!insertmacro MUI_PAGE_FINISH

; 卸载页面
!insertmacro MUI_UNPAGE_INSTFILES

; 语言文件
!insertmacro MUI_LANGUAGE "SimpChinese"

; 版本信息
VIProductVersion "{{VI_VERSION}}"
VIAddVersionKey /LANG=${LANG_SIMPCHINESE} "ProductName" "${PRODUCT_NAME}"
VIAddVersionKey /LANG=${LANG_SIMPCHINESE} "Comments" "Go版本管理工具"
VIAddVersionKey /LANG=${LANG_SIMPCHINESE} "CompanyName" "${PRODUCT_PUBLISHER}"
VIAddVersionKey /LANG=${LANG_SIMPCHINESE} "LegalTrademarks" "${PRODUCT_NAME}是${PRODUCT_PUBLISHER}的商标"
VIAddVersionKey /LANG=${LANG_SIMPCHINESE} "LegalCopyright" "© ${PRODUCT_PUBLISHER}"
VIAddVersionKey /LANG=${LANG_SIMPCHINESE} "FileDescription" "${PRODUCT_NAME} 安装程序"
VIAddVersionKey /LANG=${LANG_SIMPCHINESE} "FileVersion" "${PRODUCT_VERSION}"

; 安装目录验证函数
Function .onVerifyInstDir
  ; 检查路径长度（Windows路径限制）
  StrLen $0 "$INSTDIR"
  IntCmp $0 240 check_write check_write path_too_long
  path_too_long:
    MessageBox MB_OK|MB_ICONSTOP "安装路径太长！请选择一个较短的路径。$\n当前路径长度：$0 字符$\n最大允许：240 字符"
    Abort
  
  check_write:
  ; 检查写入权限
  ClearErrors
  CreateDirectory "$INSTDIR\__test__"
  IfErrors dir_error dir_ok
  dir_error:
    MessageBox MB_OK|MB_ICONSTOP "无法在选定路径创建目录！$\n请检查您是否有足够的权限，或选择其他安装路径。$\n$\n建议选择以下位置之一：$\n• $LOCALAPPDATA\envm$\n• $PROFILE\envm$\n• D:\envm"
    Abort
  dir_ok:
    RMDir "$INSTDIR\__test__" ; 清理测试目录
FunctionEnd

; 初始化函数
Function .onInit
  ; 设置更好的默认安装路径
  ; 优先使用用户文件夹，避免需要管理员权限
  StrCpy $INSTDIR "$LOCALAPPDATA\${PRODUCT_NAME}"
  
  ; 检查是否已经安装
  ReadRegStr $R0 ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "UninstallString"
  StrCmp $R0 "" show_welcome
  
  ; 如果已安装，获取之前的安装路径
  ReadRegStr $R1 ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "InstallLocation"
  StrCmp $R1 "" use_current_path
  StrCpy $INSTDIR "$R1"
  
  use_current_path:
  MessageBox MB_YESNOCANCEL|MB_ICONEXCLAMATION \
  "${PRODUCT_NAME} 已经安装在：$\n$R1$\n$\n点击 '是' 卸载旧版本然后重新安装$\n点击 '否' 继续安装（覆盖现有版本）$\n点击 '取消' 退出安装" \
  IDYES uninst IDNO show_welcome
  Abort
  
  uninst:
    ; 运行卸载程序
    ExecWait '$R0 _?=$R1'
    
  show_welcome:
FunctionEnd

; 安装组件
Section "核心程序" SEC01
  SectionIn RO  ; 必须安装
  ; 设置输出路径到安装目录
  SetOutPath "$INSTDIR"
  SetOverwrite ifnewer
  
  ; 复制主程序文件 (使用动态文件名)
  File /oname=envm.exe "${BINARY_FILE}"
  
  ; 创建用户目录和配置
  Call CreateUserDirectories
  
  ; 设置环境变量
  Call SetEnvironmentVariables
  
  ; 添加到系统PATH
  Call AddToPath
SectionEnd

Section "开始菜单快捷方式" SEC02
  ; 创建开始菜单快捷方式
  CreateDirectory "$SMPROGRAMS\${PRODUCT_NAME}"
  CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\${PRODUCT_NAME}.lnk" "$INSTDIR\envm.exe" "" "$INSTDIR\envm.exe" 0
  CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\卸载 ${PRODUCT_NAME}.lnk" "$INSTDIR\uninst.exe"
SectionEnd

Section "桌面快捷方式" SEC03
  ; 创建桌面快捷方式
  CreateShortCut "$DESKTOP\${PRODUCT_NAME}.lnk" "$INSTDIR\envm.exe" "" "$INSTDIR\envm.exe" 0
SectionEnd

; 组件描述
!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
  !insertmacro MUI_DESCRIPTION_TEXT ${SEC01} "envm 核心程序文件和环境配置（必须安装）"
  !insertmacro MUI_DESCRIPTION_TEXT ${SEC02} "在开始菜单创建程序快捷方式"
  !insertmacro MUI_DESCRIPTION_TEXT ${SEC03} "在桌面创建程序快捷方式"
!insertmacro MUI_FUNCTION_DESCRIPTION_END

Section -AdditionalIcons
  WriteIniStr "$INSTDIR\${PRODUCT_NAME}.url" "InternetShortcut" "URL" "${PRODUCT_WEB_SITE}"
  CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\网站.lnk" "$INSTDIR\${PRODUCT_NAME}.url"
SectionEnd

Section -Post
  WriteUninstaller "$INSTDIR\uninst.exe"
  WriteRegStr HKCU "${PRODUCT_DIR_REGKEY}" "" "$INSTDIR\envm.exe"
  WriteRegStr HKCU "${PRODUCT_UNINST_KEY}" "DisplayName" "$(^Name)"
  WriteRegStr HKCU "${PRODUCT_UNINST_KEY}" "UninstallString" "$INSTDIR\uninst.exe"
  WriteRegStr HKCU "${PRODUCT_UNINST_KEY}" "InstallLocation" "$INSTDIR"
  WriteRegStr HKCU "${PRODUCT_UNINST_KEY}" "DisplayIcon" "$INSTDIR\envm.exe"
  WriteRegStr HKCU "${PRODUCT_UNINST_KEY}" "DisplayVersion" "${PRODUCT_VERSION}"
  WriteRegStr HKCU "${PRODUCT_UNINST_KEY}" "URLInfoAbout" "${PRODUCT_WEB_SITE}"
  WriteRegStr HKCU "${PRODUCT_UNINST_KEY}" "Publisher" "${PRODUCT_PUBLISHER}"
  WriteRegDWORD HKCU "${PRODUCT_UNINST_KEY}" "NoModify" 1
  WriteRegDWORD HKCU "${PRODUCT_UNINST_KEY}" "NoRepair" 1
  
  ; 估算安装大小（约10MB）
  WriteRegDWORD HKCU "${PRODUCT_UNINST_KEY}" "EstimatedSize" 10240
SectionEnd

; 创建用户目录函数
Function CreateUserDirectories
  ; 获取用户主目录
  ReadEnvStr $0 "USERPROFILE"
  
  ; 创建 .govm 目录
  CreateDirectory "$0\.govm"
  CreateDirectory "$0\.govm\go"
  CreateDirectory "$0\.govm\versions"
  
  ; 创建默认配置文件
  FileOpen $1 "$0\.govm\settings.json" w
  FileWrite $1 '{'
  FileWrite $1 '$\n  "download_url": "https://golang.org/dl/",'
  FileWrite $1 '$\n  "mirror_url": "https://golang.google.cn/dl/",'
  FileWrite $1 '$\n  "use_mirror": false'
  FileWrite $1 '$\n}'
  FileClose $1
FunctionEnd

; 设置环境变量函数
Function SetEnvironmentVariables
  ; 获取用户主目录
  ReadEnvStr $0 "USERPROFILE"
  
  ; 设置 GOVM_HOME 环境变量
  WriteRegExpandStr HKCU "Environment" "GOVM_HOME" "$0\.govm"
  
  ; 设置 GOVM_SYMLINK 环境变量
  WriteRegExpandStr HKCU "Environment" "GOVM_SYMLINK" "$0\.govm\go"
  
  ; 通知系统环境变量已更改
  SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
FunctionEnd

; 添加到PATH函数
Function AddToPath
  ; 添加安装目录到用户PATH
  Push "$INSTDIR"
  Call AddToUserPath
  
  ; 获取用户主目录
  ReadEnvStr $0 "USERPROFILE"
  
  ; 添加 GOVM_SYMLINK\bin 到PATH
  Push "$0\.govm\go\bin"
  Call AddToUserPath
FunctionEnd

; 添加到用户PATH的辅助函数
Function AddToUserPath
  Exch $0 ; 要添加的路径
  Push $1
  Push $2
  Push $3
  
  ; 读取当前用户PATH
  ReadRegStr $1 HKCU "Environment" "PATH"
  
  ; 检查路径是否已存在
  StrLen $2 "$0"
  StrLen $3 "$1"
  
  ; 如果PATH为空，直接设置
  StrCmp $1 "" add_path
  
  ; 检查是否已包含该路径
  Push $1
  Push $0
  Call StrStr
  Pop $2
  StrCmp $2 "" add_path done
  
  add_path:
    ; 添加到PATH末尾
    StrCmp $1 "" 0 +2
    StrCpy $1 "$0"
    Goto write_path
    StrCpy $1 "$1;$0"
    
  write_path:
    WriteRegExpandStr HKCU "Environment" "PATH" "$1"
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
  
  done:
  Pop $3
  Pop $2
  Pop $1
  Pop $0
FunctionEnd

; 字符串搜索函数
Function StrStr
  Exch $R1 ; 要搜索的字符串 (needle)
  Exch 
  Exch $R2 ; 被搜索的字符串 (haystack)
  Push $R3
  Push $R4
  Push $R5
  StrLen $R3 $R1
  StrCpy $R4 0
  loop:
    StrCpy $R5 $R2 $R3 $R4
    StrCmp $R5 $R1 done
    StrCmp $R5 "" done
    IntOp $R4 $R4 + 1
    Goto loop
  done:
  StrCpy $R1 $R5
  Pop $R5
  Pop $R4
  Pop $R3
  Pop $R2
  Exch $R1
FunctionEnd

; 卸载部分
Function un.onUninstSuccess
  HideWindow
  MessageBox MB_ICONINFORMATION|MB_OK "$(^Name) 已成功从您的计算机中移除。"
FunctionEnd

Function un.onInit
  MessageBox MB_ICONQUESTION|MB_YESNO|MB_DEFBUTTON2 "您确实要完全移除 $(^Name) ，其及所有的组件？" IDYES +2
  Abort
FunctionEnd

Section Uninstall
  ; 从PATH中移除
  Call un.RemoveFromPath
  
  ; 删除环境变量
  DeleteRegValue HKCU "Environment" "GOVM_HOME"
  DeleteRegValue HKCU "Environment" "GOVM_SYMLINK"
  
  ; 删除文件
  Delete "$INSTDIR\${PRODUCT_NAME}.url"
  Delete "$INSTDIR\uninst.exe"
  Delete "$INSTDIR\envm.exe"
  
  ; 删除快捷方式
  Delete "$SMPROGRAMS\${PRODUCT_NAME}\卸载 ${PRODUCT_NAME}.lnk"
  Delete "$SMPROGRAMS\${PRODUCT_NAME}\网站.lnk"
  Delete "$SMPROGRAMS\${PRODUCT_NAME}\${PRODUCT_NAME}.lnk"
  Delete "$DESKTOP\${PRODUCT_NAME}.lnk"
  
  ; 删除目录
  RMDir "$SMPROGRAMS\${PRODUCT_NAME}"
  RMDir /r "$INSTDIR"
  
  ; 删除注册表键
  DeleteRegKey ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}"
  DeleteRegKey HKCU "${PRODUCT_DIR_REGKEY}"
  
  ; 通知系统环境变量已更改
  SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
  
  SetAutoClose true
SectionEnd

; 从PATH中移除的函数
Function un.RemoveFromPath
  ; 移除安装目录
  Push "$INSTDIR"
  Call un.RemoveFromUserPath
  
  ; 移除 GOVM_SYMLINK\bin
  ReadEnvStr $0 "USERPROFILE"
  Push "$0\.govm\go\bin"
  Call un.RemoveFromUserPath
FunctionEnd

; 从用户PATH中移除路径的辅助函数
Function un.RemoveFromUserPath
  Exch $0 ; 要移除的路径
  Push $1
  Push $2
  Push $3
  Push $4
  
  ; 读取当前PATH
  ReadRegStr $1 HKCU "Environment" "PATH"
  StrCpy $2 ""
  StrLen $3 "$0"
  
  ; 分割PATH并重建（不包含要移除的路径）
  loop:
    StrCpy $4 $1 1
    StrCmp $4 "" done
    StrCmp $4 ";" next_char
    
    ; 查找下一个分号
    Push $1
    Push ";"
    Call un.StrStr
    Pop $4
    StrCmp $4 "" last_entry
    
    ; 提取当前条目
    StrLen $4 $4
    IntOp $4 $4 - 1
    StrCpy $4 $1 $4
    
    ; 检查是否为要移除的路径
    StrCmp $4 $0 skip_entry
    StrCmp $2 "" 0 +2
    StrCpy $2 $4
    Goto +2
    StrCpy $2 "$2;$4"
    
    skip_entry:
    ; 移动到下一个条目
    Push $1
    Push ";"
    Call un.StrStr
    Pop $1
    StrCpy $1 $1 "" 1
    Goto loop
    
    last_entry:
    ; 处理最后一个条目
    StrCmp $1 $0 done
    StrCmp $2 "" 0 +2
    StrCpy $2 $1
    Goto +2
    StrCpy $2 "$2;$1"
    Goto done
    
    next_char:
    StrCpy $1 $1 "" 1
    Goto loop
  
  done:
  WriteRegExpandStr HKCU "Environment" "PATH" "$2"
  
  Pop $4
  Pop $3
  Pop $2
  Pop $1
  Pop $0
FunctionEnd

; 卸载时的字符串搜索函数
Function un.StrStr
  Exch $R1 ; 要搜索的字符串
  Exch 
  Exch $R2 ; 被搜索的字符串
  Push $R3
  Push $R4
  Push $R5
  StrLen $R3 $R1
  StrCpy $R4 0
  loop:
    StrCpy $R5 $R2 $R3 $R4
    StrCmp $R5 $R1 done
    StrCmp $R5 "" done
    IntOp $R4 $R4 + 1
    Goto loop
  done:
  StrCpy $R1 $R5
  Pop $R5
  Pop $R4
  Pop $R3
  Pop $R2
  Exch $R1
FunctionEnd