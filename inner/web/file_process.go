package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Bar struct {
	percent float32     // 百分比
	cur     int64       // 当前进度位置
	total   int64       // 总进度
	rate    string      // 进度条
	graph   string      // 显示符号
	lock    *sync.Mutex // 读写锁，正确打印
}

var bar *Bar

func NewOptionWithGraph(start, total int64, graph string) *Bar {
	bar = new(Bar)
	bar.cur = start
	bar.total = total
	if graph == "" {
		bar.graph = ">"
	}
	bar.lock = new(sync.Mutex)
	bar.percent = bar.getPercent()

	return bar
}

func (bar *Bar) getPercent() float32 {
	return float32(bar.cur) / float32(bar.total) * 100
}

func NewOption(start, total int64) *Bar {

	return NewOptionWithGraph(start, total, "")
}

func (bar *Bar) Play() {
	bar.lock.Lock()
	bar.percent = bar.getPercent()
	i := int(bar.percent)

	bar.rate = strings.Repeat(bar.graph, i)

	defer bar.lock.Unlock()
	fmt.Printf("\r[%-50s]%0.2f%% %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}

func (bar *Bar) Write(p []byte) (int, error) {
	n := len(p)
	bar.cur += int64(n)
	bar.Play()
	return n, nil
}

func Process() {
	fmt.Println("Download Started")

	fileUrl := "https://dl.google.com/go/go1.11.1.src.tar.gz"
	err := DownloadFile("go1.11.1.src.tar.gz", fileUrl)
	if err != nil {
		panic(err)
	}

	fmt.Println("Download Finished")
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(filepath string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	parseInt, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	// Create our progress reporter and pass it to be used alongside our writer
	counter := NewOption(0, parseInt)
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	out.Close()
	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")
	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}
