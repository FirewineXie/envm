package web_java

import (
	"fmt"
	"github.com/FirewineXie/envm/util"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	DefaultURL = "https://www.oracle.com/cn/java/technologies/downloads/archive/"
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

// NewCollector 返回采集器实例
func NewCollector(url string) (*Collector, error) {
	if url == "" {
		url = DefaultURL
	}
	c := Collector{
		url: url,
	}
	resp, err := http.Get(c.url)
	if err != nil {
		return nil, err
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

// LatestSubPackage Compressed Archive
// 找到第一个就返回
func (c *Collector) LatestSubPackage(goos, goarch string) (p *util.Package, err error) {
	var packageUrl string
	var sha256 string
	title := strings.Title(goos)
	c.doc.Find("table.otable-w2").Each(func(i int, div *goquery.Selection) {
		div.Find("tr").Find("td").Each(func(i int, selection *goquery.Selection) {

			if packageUrl != "" {
				return
			}
			product := selection.Text()
			// 只下载压缩文件
			//cases.Title(goos,cases.Option())

			if strings.Contains(product, "Compressed Archive") && strings.Contains(product, title) {
				nodes := selection.Find("a").Nodes
				if len(nodes) > 1 {

					newPackageUrl := nodes[0].Attr[0].Val
					newSha256 := nodes[1].Attr[0].Val
					// 判断版本是否一样
					if strings.Contains(packageUrl, goarch) {
						packageUrl = newPackageUrl
						sha256 = newSha256
					}
				}
			}

		})

	})
	p = &util.Package{}
	p.URL = packageUrl
	p.Algorithm = "sha256"
	p.Checksum = sha256
	p.OS = goos
	p.Arch = goarch
	return

}

//func (c *Collector) findPackages(table *goquery.Selection) (pkgs []*Package) {
//	table.Find("tr").Find("td").Each(func(i int, selection *goquery.Selection) {
//		nodes := selection.Find("a").Nodes
//		if len(nodes) > 1 {
//			packgeUrl := nodes[0].Attr[0].Val
//			sha256 := nodes[1].Attr[0].Val
//			pkgs = append(pkgs, &Package{
//				FileName:  td.Eq(0).Find("a").Text(),
//				URL:       nodes[0].Attr[0].Val,
//				Kind:      td.Eq(1).Text(),
//				OS:        td.Eq(2).Text(),
//				Arch:      td.Eq(3).Text(),
//				Size:      td.Eq(4).Text(),
//				Checksum:  td.Eq(5).Text(),
//				Algorithm: "sha256",
//			})
//
//		}
//
//	})
//	return pkgs
//}

// LatestFiveVersion 返回最新的5个大版本
func (c *Collector) LatestFiveVersion() (items []*util.Version, err error) {
	var stopInt int
	c.doc.Find("ul.icn-ulist").Find("li.icn-chevron-right").Each(func(i int, div *goquery.Selection) {
		versionDescribe := div.Find("a").Text()
		val, ok := div.Find("a").Attr("href")
		if !ok {
			return
		}
		version := strings.TrimPrefix(versionDescribe, "Java SE ")
		version = strings.Replace(version, "\n", " ", -1)
		if version == "7" {
			stopInt = i
			return
		}
		if stopInt != 0 {
			return
		}
		items = append(items, &util.Version{
			Name: version,
			Packages: []*util.Package{
				{
					URL: "https://www.oracle.com/" + val,
				},
			},
		})
	})
	return items, nil
}
