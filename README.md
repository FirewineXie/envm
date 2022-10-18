# 写一个类似于nvm 的go 开源工具

> 基于GoLand 的切换环境的功能

> 在go module 的基础上进行切换，没有1.13 下面的不进行考虑

## 安装过程

1. 在电脑自己目录下面新建 .govm 文件夹
2. 在系统变量设置
    1. `GOVM_HOME`  example :  C:\Users\username\.govm
    2. `GOVM_SYMLINK` example : C:\Users\username\.govm\go

3. 尝试运行govm 是否可以正常运行 example: govm arch
4. 在`GOVM_HOME`里面修改settings配置文件，
    1. 暂时只支持修改下载目录

## 尾注

感谢 `gvm`,`nvm` 提供的灵感和代码的实现