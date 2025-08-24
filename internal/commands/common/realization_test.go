package common

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_getCurrentVersion(t *testing.T) {
	Convey("验证go 当前版本", t, func() {
		version := GetCurrentVersion("go")
		So(version, ShouldNotBeEmpty)
	})
}

func Test_getInstalled(t *testing.T) {
	Convey("测试已经安装的go 版本", t, func() {
		versions := GetInstalled("D:\\programs", "go")
		fmt.Println(versions)
		So(versions, ShouldNotBeNil)
	})
}
