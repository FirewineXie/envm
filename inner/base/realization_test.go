package base

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_getCurrentVersion(t *testing.T) {
	Convey("验证go 当前版本", t, func() {
		version := getCurrentVersion()
		So(version, ShouldEqual, "go1.17.1")
	})
}

func Test_getInstalled(t *testing.T) {
	Convey("测试已经安装的go 版本", t, func() {
		fmt.Println(getInstalled("D:\\ProgramData\\Go"))
		
	})
}
