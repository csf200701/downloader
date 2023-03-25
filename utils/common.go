package utils

import (
	"downloader/config"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func CalcSize(size int64) string {
	var show float64
	var unit string
	if size < 1024 {
		show = float64(size)
		unit = "B"
	} else if size < int64(math.Pow(1024, 2)) {
		show = float64(size) / math.Pow(1024, 1)
		unit = "KB"
	} else if size < int64(math.Pow(1024, 3)) {
		show = float64(size) / math.Pow(1024, 2)
		unit = "MB"
	} else if size < int64(math.Pow(1024, 4)) { //千兆
		show = float64(size) / math.Pow(1024, 3)
		unit = "GB"
	} else if size < int64(math.Pow(1024, 5)) { //太字节
		show = float64(size) / math.Pow(1024, 4)
		unit = "TB"
	} else if size < int64(math.Pow(1024, 6)) { //拍字节
		show = float64(size) / math.Pow(1024, 5)
		unit = "PB"
	} else if size < int64(math.Pow(1024, 7)) { //艾字节
		show = float64(size) / math.Pow(1024, 6)
		unit = "EB"
	} else if size < int64(math.Pow(1024, 8)) { //泽字节
		show = float64(size) / math.Pow(1024, 7)
		unit = "ZB"
	} else if size < int64(math.Pow(1024, 9)) { //尧字节
		show = float64(size) / math.Pow(1024, 8)
		unit = "YB"
	} else if size < int64(math.Pow(1024, 10)) { //
		show = float64(size) / math.Pow(1024, 9)
		unit = "BB"
	}
	return strconv.FormatFloat(show, 'f', 2, 64) + unit
}

func CalcProcess(dm int64) int {
	var process int
	if dm < 40 {
		process = 1
	} else if dm < 1024 {
		process = 5
	} else {
		process = 10
	}
	return process
}

func ComponentUrl(componentName string, componentVersion string) (string, *config.Component, error) {
	if len(componentName) == 0 {
		return "", nil, errors.New("组件不能为空")
	}
	var url string
	var base string
	var comp *config.Component
	c := config.C
	for _, component := range c.Components {
		if component.Name == componentName {
			if len(componentVersion) > 0 {
				for _, version := range component.Versions {
					if version.Name == componentVersion {
						url = version.Url
						base = component.Base
						comp = &component
						break
					}
				}
			} else {
				url = component.Versions[0].Url
				base = component.Base
				comp = &component
			}
			break
		}
	}
	if len(url) == 0 {
		return "", nil, errors.New("获取不到组件所对应的URL地址")
	}
	if len(base) > 0 && strings.Index(url, strings.ToLower("https|http")) == -1 {
		if strings.HasSuffix(base, "/") {
			base = strings.TrimSuffix(base, "/")
		}
		if strings.HasSuffix(url, "/") {
			url = strings.TrimSuffix(url, "/")
		}
		url = base + "/" + url
	}

	return url, comp, nil
}

func PorxyClient(proxyHost, proxyUsername, proxyPassword string, timeout time.Duration) (*http.Client, string) {
	var auth string
	if len(proxyUsername) > 0 {
		auth += "username=" + proxyUsername
	}
	if len(proxyPassword) > 0 {
		if len(auth) > 0 {
			auth += "&"
		}
		auth += "password=" + proxyPassword
	}
	//HTTP代理
	proxy := fmt.Sprintf("http://%s", proxyHost)
	proxyAddress, _ := url.Parse(proxy)
	getClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyAddress),
		},
	}
	return getClient, auth
}
