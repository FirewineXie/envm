package web_go

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/FirewineXie/envm/util"
	"io/ioutil"
	"testing"

	"github.com/PuerkitoBio/goquery"

	. "github.com/smartystreets/goconvey/convey"
)

func getCollector() (*Collector, error) {
	b, err := ioutil.ReadFile("./testdata/goDownload.html")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return &Collector{
		url: DefaultURL,
		doc: doc,
	}, nil
}

func Test_findPackages(t *testing.T) {
	Convey("查找目标go版本下的安装包列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		pkgs := c.findPackages(c.doc.Find("#stable").Next().Find("table").First())
		//So(len(pkgs), ShouldEqual, 15)
		So(pkgs[1].Algorithm, ShouldEqual, "SHA256")
		So(pkgs[1].FileName, ShouldEqual, "go1.12.4.darwin-amd64.tar.gz")
		So(pkgs[1].Kind, ShouldEqual, util.ArchiveKind)
		So(pkgs[1].OS, ShouldEqual, "macOS")
		So(pkgs[1].Arch, ShouldEqual, "x86-64")
		So(pkgs[1].Size, ShouldEqual, "122MB")
		So(pkgs[1].Checksum, ShouldEqual, "50af1aa6bf783358d68e125c5a72a1ba41fb83cee8f25b58ce59138896730a49")
	})
}

func TestStableVersions(t *testing.T) {
	Convey("查询stable状态的go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.StableVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldEqual, 2)
		So(items[0].Name, ShouldEqual, "1.12.4")
		So(len(items[0].Packages), ShouldEqual, 15)
		So(items[1].Name, ShouldEqual, "1.11.9")
		So(len(items[1].Packages), ShouldEqual, 15)
	})
}

func TestArchivedVersions(t *testing.T) {
	Convey("查询archived状态的go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.ArchivedVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldEqual, 64)

		So(items[0].Name, ShouldEqual, "1.12.3")
		So(len(items[0].Packages), ShouldEqual, 15)
	})
}

func TestAllVersions(t *testing.T) {
	Convey("查询所有go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.AllVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldEqual, 66)
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

func TestVersionGO_FindPackage(t *testing.T) {
	Convey("测试FindPackage架构映射功能", t, func() {
		// 创建测试版本数据
		version := &VersionGO{
			Version: util.Version{
				Name: "1.21.0",
				Packages: []*util.Package{
					{
						FileName: "go1.21.0.linux-amd64.tar.gz",
						Kind:     util.ArchiveKind,
						OS:       "Linux",
						Arch:     "x86-64",
					},
					{
						FileName: "go1.21.0.linux-386.tar.gz", 
						Kind:     util.ArchiveKind,
						OS:       "Linux",
						Arch:     "x86",
					},
					{
						FileName: "go1.21.0.linux-arm64.tar.gz",
						Kind:     util.ArchiveKind,
						OS:       "Linux", 
						Arch:     "ARMv8",
					},
					{
						FileName: "go1.21.0.windows-amd64.zip",
						Kind:     util.ArchiveKind,
						OS:       "Windows",
						Arch:     "x86-64",
					},
				},
			},
		}

		Convey("x86_64应该找到amd64包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "x86_64")
			So(err, ShouldBeNil)
			So(pkg, ShouldNotBeNil)
			So(pkg.FileName, ShouldEqual, "go1.21.0.linux-amd64.tar.gz")
		})

		Convey("amd64应该找到amd64包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "amd64")
			So(err, ShouldBeNil)
			So(pkg, ShouldNotBeNil)
			So(pkg.FileName, ShouldEqual, "go1.21.0.linux-amd64.tar.gz")
		})

		Convey("i386应该找到386包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "i386")
			So(err, ShouldBeNil)
			So(pkg, ShouldNotBeNil)
			So(pkg.FileName, ShouldEqual, "go1.21.0.linux-386.tar.gz")
		})

		Convey("386应该找到386包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "386")
			So(err, ShouldBeNil)
			So(pkg, ShouldNotBeNil)
			So(pkg.FileName, ShouldEqual, "go1.21.0.linux-386.tar.gz")
		})

		Convey("aarch64应该找到arm64包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "aarch64")
			So(err, ShouldBeNil)
			So(pkg, ShouldNotBeNil)
			So(pkg.FileName, ShouldEqual, "go1.21.0.linux-arm64.tar.gz")
		})

		Convey("不存在的架构应该返回错误", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "linux", "nonexistent")
			So(err, ShouldEqual, util.ErrPackageNotFound)
			So(pkg, ShouldBeNil)
		})

		Convey("Windows x86_64应该找到对应包", func() {
			pkg, err := version.FindPackage(util.ArchiveKind, "windows", "x86_64")
			So(err, ShouldBeNil)
			So(pkg, ShouldNotBeNil)
			So(pkg.FileName, ShouldEqual, "go1.21.0.windows-amd64.zip")
		})
	})
}
