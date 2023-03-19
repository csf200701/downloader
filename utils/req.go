package utils

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"time"
)

type RemoteFileInfo struct {
	Total    int64
	FileName string
}
type Request struct {
	url string
}

func NewRequest(url string) *Request {
	return &Request{url}
}

func (r *Request) Total() (*RemoteFileInfo, error) { //Content-Disposition: attachment; filename=Postman-win64-Setup.exe
	totalReq, err := http.NewRequest("HEAD", r.url, nil)
	if err != nil {
		return nil, err
	}
	totalClient := http.Client{Timeout: time.Second * 10}
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
	getReq, err := http.NewRequest("GET", r.url, nil)
	getReq.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", 0, 0))
	getClient := http.Client{Timeout: time.Second * 60 * 10}
	getRsp, err := getClient.Do(getReq)
	if err != nil {
		return 0, err
	}
	return getRsp.ContentLength, nil
}

func (r *Request) Content(start, end int64) (io.ReadCloser, int64, error) {
	getReq, err := http.NewRequest("GET", r.url, nil)
	getReq.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	getClient := http.Client{Timeout: time.Second * 60 * 10}
	getRsp, err := getClient.Do(getReq)
	if err != nil {
		return nil, 0, err
	}
	// fmt.Println(start, end, getRsp)
	return getRsp.Body, getRsp.ContentLength, nil
}
