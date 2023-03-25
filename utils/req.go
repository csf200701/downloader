package utils

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RemoteFileInfo struct {
	Total    int64
	FileName string
}
type Request struct {
	url   string
	proxy *Proxy
}

type Proxy struct {
	host     string
	username string
	password string
}

func NewProxy(host, username, password string) *Proxy {
	return &Proxy{host, username, password}
}

func NewRequest(url string) *Request {
	return &Request{url: url}
}

func NewProxyRequest(url string, proxy *Proxy) *Request {
	return &Request{url: url, proxy: proxy}
}

func (r *Request) Total() (*RemoteFileInfo, error) {
	var totalClient *http.Client
	var totalReq *http.Request
	var err error
	if r.proxy == nil {
		totalClient = &http.Client{Timeout: time.Second * 10}
		totalReq, err = http.NewRequest("HEAD", r.url, nil)
		if err != nil {
			return nil, err
		}
	} else {
		var auth string
		totalClient, auth = PorxyClient(r.proxy.host, r.proxy.username, r.proxy.password, time.Second*60)
		totalReq, err = http.NewRequest("HEAD", r.url, strings.NewReader(auth))
		if err != nil {
			return nil, err
		}
	}

	totalRsp, err := totalClient.Do(totalReq)
	if err != nil {
		return nil, err
	}
	defer totalRsp.Body.Close()
	totalSize, err := strconv.ParseInt(totalRsp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return nil, err
	}
	remoteFileInfo := &RemoteFileInfo{Total: totalSize}
	contentDisposition := totalRsp.Header.Get("Content-Disposition")
	if len(contentDisposition) > 0 {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err == nil {
			filename := params["filename"]
			remoteFileInfo.FileName = filename
		}
	}

	return remoteFileInfo, nil
}

func (r *Request) Total_Get() (int64, error) {
	var getClient *http.Client
	var getReq *http.Request
	var err error
	if r.proxy == nil {
		getClient = &http.Client{Timeout: time.Second * 60 * 10}
		getReq, err = http.NewRequest("GET", r.url, nil)
	} else {
		var auth string
		getClient, auth = PorxyClient(r.proxy.host, r.proxy.username, r.proxy.password, time.Second*60*10)
		getReq, err = http.NewRequest("GET", r.url, strings.NewReader(auth))
	}
	getReq.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", 0, 0))
	getRsp, err := getClient.Do(getReq)
	if err != nil {
		return 0, err
	}

	return getRsp.ContentLength, nil
}

func (r *Request) Content(start, end int64) (io.ReadCloser, int64, error) {
	var getClient *http.Client
	var getReq *http.Request
	var err error
	if r.proxy == nil {
		getClient = &http.Client{Timeout: time.Second * 60 * 10}
		getReq, err = http.NewRequest("GET", r.url, nil)
	} else {
		var auth string
		getClient, auth = PorxyClient(r.proxy.host, r.proxy.username, r.proxy.password, time.Second*60*10)
		getReq, err = http.NewRequest("GET", r.url, strings.NewReader(auth))

	}

	getReq.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	getRsp, err := getClient.Do(getReq)
	if err != nil {
		return nil, 0, err
	}

	return getRsp.Body, getRsp.ContentLength, nil
}
