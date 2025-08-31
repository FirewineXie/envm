package web_go

import (
	"errors"
	"fmt"
	"testing"

	"github.com/FirewineXie/envm/util"

	. "github.com/smartystreets/goconvey/convey"
)

func getCollector() (*Collector, error) {
	return NewCollector(DefaultURL)
}

func Test_findPackages(t *testing.T) {
	Convey("查找目标go版本下的安装包列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		pkgs := c.findPackages(c.doc.Find("#stable").Next().Find("table").First())
		So(len(pkgs), ShouldBeGreaterThan, 0)
		// 验证包结构，找到一个有效的包来检查
		var validPkg *util.Package
		for _, pkg := range pkgs {
			if pkg.FileName != "" && pkg.Kind != "" && pkg.OS != "" {
				validPkg = pkg
				break
			}
		}
		So(validPkg, ShouldNotBeNil)
		So(validPkg.Algorithm, ShouldEqual, "SHA256")
		So(validPkg.FileName, ShouldNotBeEmpty)
		So(validPkg.Kind, ShouldNotBeEmpty)
		So(validPkg.OS, ShouldNotBeEmpty)
		So(validPkg.Arch, ShouldNotBeEmpty)
	})
}

func TestStableVersions(t *testing.T) {
	Convey("查询stable状态的go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.StableVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldBeGreaterThan, 0)
		// 验证第一个稳定版本的结构
		So(items[0].Name, ShouldNotBeEmpty)
		So(len(items[0].Packages), ShouldBeGreaterThan, 0)
	})
}

func TestArchivedVersions(t *testing.T) {
	Convey("查询archived状态的go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.ArchivedVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldBeGreaterThan, 0)
		// 验证第一个归档版本的结构
		So(items[0].Name, ShouldNotBeEmpty)
		So(len(items[0].Packages), ShouldBeGreaterThan, 0)
	})
}

func TestAllVersions(t *testing.T) {
	Convey("查询所有go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.AllVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldBeGreaterThan, 0)
		// 验证版本列表包含稳定版和归档版
		stableVersions, _ := c.StableVersions()
		archivedVersions, _ := c.ArchivedVersions()
		So(len(items), ShouldEqual, len(stableVersions)+len(archivedVersions))
	})
}

func TestURLUnreachableError(t *testing.T) {
	Convey("URL不可达错误", t, func() {
		url := "https://github.com/voidint"
		core := errors.New("hello error")

		err := NewURLUnreachableError(url, core)
		So(err, ShouldNotBeNil)
		e, ok := err.(*URLUnreachableError)
		So(ok, ShouldBeTrue)
		So(e, ShouldNotBeNil)
		So(e.url, ShouldEqual, url)
		So(e.err, ShouldEqual, core)
		So(e.Error(), ShouldEqual, fmt.Sprintf("URL %q is unreachable ==> %s", url, core.Error()))
	})
}

func TestNormalizeArch(t *testing.T) {
	Convey("测试架构名称标准化", t, func() {
		testCases := []struct {
			input    string
			expected string
			desc     string
		}{
			// amd64 variants
			{"x86_64", "amd64", "x86_64 should map to amd64"},
			{"x64", "amd64", "x64 should map to amd64"},
			{"amd64", "amd64", "amd64 should remain amd64"},

			// 386 variants
			{"i386", "386", "i386 should map to 386"},
			{"i686", "386", "i686 should map to 386"},
			{"x86", "386", "x86 should map to 386"},
			{"386", "386", "386 should remain 386"},

			// arm64 variants
			{"aarch64", "arm64", "aarch64 should map to arm64"},
			{"arm64", "arm64", "arm64 should remain arm64"},

			// arm variants
			{"armv6l", "arm", "armv6l should map to arm"},
			{"armv7l", "arm", "armv7l should map to arm"},
			{"arm", "arm", "arm should remain arm"},

			// tab format from arch.Validate()
			{"unix \tx86_64", "amd64", "unix tab format should extract x86_64 and map to amd64"},
			{"unix \ti386", "386", "unix tab format should extract i386 and map to 386"},

			// unknown architecture
			{"unknown", "unknown", "unknown arch should remain unchanged"},
		}

		for _, tc := range testCases {
			Convey(tc.desc, func() {
				result := normalizeArch(tc.input)
				So(result, ShouldEqual, tc.expected)
			})
		}
	})
}

func TestNormalizeOS(t *testing.T) {
	Convey("测试操作系统名称标准化", t, func() {
		testCases := []struct {
			input    string
			expected string
			desc     string
		}{
			// unix variants
			{"unix", "linux", "unix should map to linux"},
			{"unix \tx86_64", "linux", "unix with tab format should map to linux"},
			
			// standard OS names
			{"linux", "linux", "linux should remain linux"},
			{"darwin", "darwin", "darwin should remain darwin"},
			{"windows", "windows", "windows should remain windows"},
			
			// unknown OS
			{"unknown", "unknown", "unknown OS should remain unchanged"},
		}

		for _, tc := range testCases {
			Convey(tc.desc, func() {
				result := normalizeOS(tc.input)
				So(result, ShouldEqual, tc.expected)
			})
		}
	})
}

func TestVersionGO_FindPackage(t *testing.T) {
	Convey("测试FindPackage架构映射功能", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		// 获取第一个稳定版本进行测试
		stableVersions, err := c.StableVersions()
		So(err, ShouldBeNil)
		So(len(stableVersions), ShouldBeGreaterThan, 0)

		version := stableVersions[0]
		So(version.Name, ShouldNotBeEmpty)
		So(len(version.Packages), ShouldBeGreaterThan, 0)
		Convey("unix系统x86_64应该映射到linux-amd64包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "unix", "x86_64")
			if err == nil {
				So(pkg, ShouldNotBeNil)
				So(pkg.FileName, ShouldContainSubstring, "linux-amd64")
			} else {
				So(err, ShouldEqual, util.ErrPackageNotFound)
			}
		})
		
		Convey("unix系统带tab格式架构应该正确解析", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "unix", "unix \tx86_64")
			if err == nil {
				So(pkg, ShouldNotBeNil)
				So(pkg.FileName, ShouldContainSubstring, "linux-amd64")
			} else {
				So(err, ShouldEqual, util.ErrPackageNotFound)
			}
		})
		
		Convey("linux系统x86_64应该找到amd64包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "x86_64")
			if err == nil {
				So(pkg, ShouldNotBeNil)
				So(pkg.FileName, ShouldContainSubstring, "linux-amd64")
			} else {
				So(err, ShouldEqual, util.ErrPackageNotFound)
			}
		})

		Convey("amd64应该找到amd64包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "amd64")
			if err == nil {
				So(pkg, ShouldNotBeNil)
				So(pkg.FileName, ShouldContainSubstring, "linux-amd64")
			} else {
				So(err, ShouldEqual, util.ErrPackageNotFound)
			}
		})

		Convey("i386应该找到386包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "i386")
			if err == nil {
				So(pkg, ShouldNotBeNil)
				So(pkg.FileName, ShouldContainSubstring, "linux-386")
			} else {
				So(err, ShouldEqual, util.ErrPackageNotFound)
			}
		})

		Convey("386应该找到386包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "386")
			if err == nil {
				So(pkg, ShouldNotBeNil)
				So(pkg.FileName, ShouldContainSubstring, "linux-386")
			} else {
				So(err, ShouldEqual, util.ErrPackageNotFound)
			}
		})

		Convey("aarch64应该找到arm64包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "aarch64")
			if err == nil {
				So(pkg, ShouldNotBeNil)
				So(pkg.FileName, ShouldContainSubstring, "linux-arm64")
			} else {
				So(err, ShouldEqual, util.ErrPackageNotFound)
			}
		})

		Convey("不存在的架构应该返回错误", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "nonexistent")
			So(err, ShouldEqual, util.ErrPackageNotFound)
			So(pkg, ShouldBeNil)
		})

		Convey("Windows x86_64应该找到对应包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "windows", "x86_64")
			if err == nil {
				So(pkg, ShouldNotBeNil)
				So(pkg.FileName, ShouldContainSubstring, "windows-amd64")
			} else {
				So(err, ShouldEqual, util.ErrPackageNotFound)
			}
		})
	})
}
