# 写一个类似于nvm 的go 开源工具

> 基于GoLand 的切换环境的功能

> 在go module 的基础上进行切换，没有1.13 下面的不进行考虑

## 实现功能

### 所实现的命令

1. govm arch 查看自己电脑 bit mode
2. govm active 默认激活环境
3. govm install 下载 版本
4. govm list 查看可用版本
8. govm uninstall 卸载 （但不能卸载当前使用环境）
9. govm root 设置 goroot 父级目录
10. govm version 版本

### 配置文件所属

#### 公有配置文件

- windows
  `C:\Users\xyjwo\AppData\Roaming\go\env`

## 尾注

感谢 `gvm`,`nvm` 提供的灵感和代码的实现