package web_node

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
)

/*
 * @Author: Firewine
 * @File: web
 * @Version: 1.0.0
 * @Date: 2024-04-05 17:21
 * @Description:
 */
var client = &http.Client{}

func SetProxy(p string, verifyssl bool) {
	if p != "" && p != "none" {
		proxyUrl, _ := url.Parse(p)
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl), TLSClientConfig: &tls.Config{InsecureSkipVerify: verifyssl}}}
	} else {
		client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: verifyssl}}}
	}
}

func DownloadContent(url string) (content []byte, err error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)

}
