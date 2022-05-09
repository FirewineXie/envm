# 写一个类似于nvvm 的go 开源工具

> 基于GoLand 的切换环境的功能

> 在go module 的基础上进行切换，没有1.13 下面的不进行考虑

## 实现功能



### 所实现的命令



1. govm arch   查看自己电脑 bit mode 
2. govm current  默认激活环境
3. govm install 下载 版本
4. govm list  查看可用版本
5. govm on  启用这个东东
6. govm off  关闭这个东东
7. govm proxy 设置proxy 进行下载
8. govm uninstall 卸载 （但不能卸载当前使用环境）
9. govm root  设置  goroot 父级目录
10. govm version 版本



### 配置文件所属



#### 公有配置文件

`C:\Users\xyjwo\AppData\Roaming\go\env`



#### 是否需要私有配置文件？

1. 如果需要那么私有配置文件，需要配置什么

   > 模块化是否打开？
