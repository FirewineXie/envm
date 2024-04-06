package web_node

import (
	"github.com/FirewineXie/envm/internal/arch"
	"github.com/FirewineXie/envm/util"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAvailable(t *testing.T) {
	Convey("查找目标node版本下的安装包列表", t, func() {
		all, lts, current, stable, unstable, npm, _ := GetAvailable()

		// all
		So(all[0], ShouldEqual, "21.7.2")
		So(lts[0], ShouldEqual, "20.12.1")
		So(current[0], ShouldEqual, "21.7.2")
		So(stable[0], ShouldEqual, "0.12.18")
		So(unstable[0], ShouldEqual, "0.11.16")
		So(npm[all[0]], ShouldEqual, "10.5.0")

	})
}

func TestVersionNode_FindPackage(t *testing.T) {
	Convey("查找目标版本", t, func() {
		GetAvailable()

		element, ok := GetMeta()["21.7.2"]
		if !ok {
			return
		}
		findPackage, err := element.FindPackage(util.ArchiveKind, runtime.GOOS, arch.Validate())

		if err != nil {
			panic(err)
		}

		t.Log(findPackage)

	})
}
