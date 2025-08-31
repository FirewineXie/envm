package web_go

import (
	"fmt"
	"github.com/FirewineXie/envm/util"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	// DefaultURL 提供go版本信息的默认网址
	DefaultURL = "https://golang.org/dl/"
	// BackupURL 备用网址（中国镜像）
	BackupURL = "https://golang.google.cn/dl/"
	// ProxyURL 代理网址（通过代理访问）
	ProxyURL = "https://goproxy.cn/golang.org/dl/"
)

// URLUnreachableError URL不可达错误
type URLUnreachableError struct {
	err error
	url string
}

// NewURLUnreachableError 返回URL不可达错误实例
func NewURLUnreachableError(url string, err error) error {
	return &URLUnreachableError{
		err: err,
		url: url,
	}
}

func (e *URLUnreachableError) Error() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("URL %q is unreachable", e.url))
	if e.err != nil {
		buf.WriteString(" ==> " + e.err.Error())
	}
	return buf.String()
}

type Collector struct {
	url string
	doc *goquery.Document
}

// GetURL 返回当前使用的URL
func (c *Collector) GetURL() string {
	return c.url
}

// NewCollector 返回采集器实例，支持URL自动回退
func NewCollector(url string) (*Collector, error) {
	if url == "" {
		url = DefaultURL
	}
	
	c := Collector{
		url: url,
	}
	
	// 尝试连接指定URL
	resp, err := http.Get(c.url)
	if err != nil {
		// 如果是默认URL失败，尝试备用URL
		if url == DefaultURL {
			fmt.Printf("主URL连接失败，尝试备用URL: %s\n", BackupURL)
			c.url = BackupURL
			resp, err = http.Get(c.url)
			if err != nil {
				return nil, fmt.Errorf("主URL和备用URL都无法连接: %v", err)
			}
		} else {
			return nil, err
		}
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, NewURLUnreachableError(c.url, nil)
	}
	c.doc, err = goquery.NewDocumentFromReader(resp.Body)
	return &c, nil
}

func (c *Collector) loadDocument() (err error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return NewURLUnreachableError(c.url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return NewURLUnreachableError(c.url, nil)
	}
	c.doc, err = goquery.NewDocumentFromReader(resp.Body)
	return err
}

func (c *Collector) findPackages(table *goquery.Selection) (pkgs []*util.Package) {
	alg := strings.TrimSuffix(table.Find("thead").Find("th").Last().Text(), " Checksum")

	table.Find("tr").Not(".first").Each(func(j int, tr *goquery.Selection) {
		td := tr.Find("td")
		pkgs = append(pkgs, &util.Package{
			FileName:  td.Eq(0).Find("a").Text(),
			URL:       td.Eq(0).Find("a").AttrOr("href", ""),
			Kind:      td.Eq(1).Text(),
			OS:        td.Eq(2).Text(),
			Arch:      td.Eq(3).Text(),
			Size:      td.Eq(4).Text(),
			Checksum:  td.Eq(5).Text(),
			Algorithm: alg,
		})
	})
	return pkgs
}

// StableVersions 返回所有稳定版本
func (c *Collector) StableVersions() (items []*VersionGO, err error) {
	c.doc.Find("#stable").NextUntil("#archive").Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}

		versionGO := &VersionGO{}
		versionGO.Name = strings.TrimPrefix(vname, "go")
		versionGO.Packages = c.findPackages(div.Find("table").First())
		items = append(items, versionGO)
	})
	return items, nil
}

// ArchivedVersions 返回已归档版本
func (c *Collector) ArchivedVersions() (items []*VersionGO, err error) {
	c.doc.Find("#archive").Find("div.toggle").Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		versionGo := &VersionGO{}
		versionGo.Name = strings.TrimPrefix(vname, "go")
		versionGo.Packages = c.findPackages(div.Find("table").First())
		items = append(items, versionGo)
	})
	return items, nil
}

// AllVersions 返回所有已知版本
func (c *Collector) AllVersions() (items []*VersionGO, err error) {
	items, err = c.StableVersions()
	if err != nil {
		return nil, err
	}
	archives, err := c.ArchivedVersions()
	if err != nil {
		return nil, err
	}
	items = append(items, archives...)
	return items, nil
}

type VersionGO struct {
	util.Version
}

// FindPackage 返回指定操作系统和硬件架构的版本包
func (v *VersionGO) FindPackage(kind, goos, goarch string) (*util.Package, error) {
	// 标准化操作系统名称映射
	goos = normalizeOS(goos)
	// 标准化架构名称映射
	goarch = normalizeArch(goarch)
	prefix := fmt.Sprintf("go%s.%s-%s", v.Name, goos, goarch)
	for i := range v.Packages {
		if v.Packages[i] == nil || !strings.EqualFold(v.Packages[i].Kind, kind) || !strings.HasPrefix(v.Packages[i].FileName, prefix) {
			continue
		}
		return v.Packages[i], nil
	}
	return nil, util.ErrPackageNotFound
}

// normalizeOS 标准化操作系统名称到Go官方命名
func normalizeOS(os string) string {
	// 处理arch.Validate()返回的格式 "unix \t架构"
	if strings.HasPrefix(os, "unix") {
		return "linux"
	}
	switch os {
	case "darwin":
		return "darwin"
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	default:
		return os
	}
}

// normalizeArch 标准化架构名称到Go官方命名
func normalizeArch(arch string) string {
	// 处理arch.Validate()返回的格式 "unix \t架构"
	if strings.Contains(arch, "\t") {
		parts := strings.Split(arch, "\t")
		if len(parts) > 1 {
			arch = strings.TrimSpace(parts[1])
		}
	}
	
	switch arch {
	case "x86_64", "x64", "amd64":
		return "amd64"
	case "i386", "i686", "x86", "386":
		return "386"
	case "aarch64", "arm64":
		return "arm64"
	case "armv6l", "armv7l", "arm":
		return "arm"
	default:
		return arch
	}
}
