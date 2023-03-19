package delete

import (
	"downloader/utils"
	"fmt"
	netUrl "net/url"
	"os"
	"strconv"
	"strings"
)

type Deleter struct {
	component string
	version   string
	process   int
}

func NewDeleter(component string, version string, process int) *Deleter {
	return &Deleter{component: component, version: version, process: process}
}

func (d *Deleter) Delete() {
	// c := config.C
	// var components []config.Component = make([]config.Component, 0)
	// for _, component := range c.Components {
	// 	if len(d.component) > 0 {
	// 		if d.component == component.Name {
	// 			components = append(components, component)
	// 			if len(d.version) > 0 {
	// 				component.Versions = make([]config.Version, 0)
	// 				for _, version := range component.Versions {
	// 					if d.version == version.Name {
	// 						component.Versions = append(component.Versions, version)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	} else {
	// 		components = append(components, component)
	// 	}
	// }

	// for _, component := range components {
	// 	var url string
	// 	var base string
	// 	for _, version := range component.Versions {
	// 		url = version.Url
	// 		base = component.Base
	// 		if len(url) == 0 {
	// 			continue
	// 		}
	// 		if len(base) > 0 && strings.Index(url, strings.ToLower("https|http")) == -1 {
	// 			if strings.HasSuffix(base, "/") {
	// 				base = strings.TrimSuffix(base, "/")
	// 			}
	// 			if strings.HasSuffix(url, "/") {
	// 				url = strings.TrimSuffix(url, "/")
	// 			}
	// 			url = base + "/" + url
	// 		}

	// 		req := utils.NewRequest(url)

	// 		remoteFile, err := req.Total()
	// 		if err != nil {
	// 			fmt.Println(component.Name+"-"+version.Name, "发送错误：", err)
	// 			return
	// 		}
	// 		total := remoteFile.Total
	// 		var process int = d.process
	// 		if process == 0 {
	// 			dm := total / 1024 / 1024
	// 			if dm < 40 {
	// 				process = 1
	// 			} else if dm < 1024 {
	// 				process = 5
	// 			} else {
	// 				process = 10
	// 			}
	// 		}
	// 		var partitionTotal = total / int64(process)
	// 		var fileName string
	// 		if len(remoteFile.FileName) > 0 {
	// 			fileName = remoteFile.FileName
	// 		} else {
	// 			urlObj, _ := netUrl.Parse(url)
	// 			urlPath := urlObj.Path
	// 			lastIdx := strings.LastIndex(urlPath, "/")
	// 			fileName = string([]rune(urlPath)[lastIdx+1:])
	// 		}

	// 		for i := 0; i < process; i++ {
	// 			start := partitionTotal * int64(i)
	// 			end := partitionTotal*(int64(i)+1) - 1
	// 			if process-1 == i {
	// 				end = total
	// 			}
	// 			prepend := utils.Short(url + ":" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10))
	// 			err = os.Remove(prepend)
	// 			if err != nil {
	// 				fmt.Println("组件："+component.Name+"，版本："+version.Name+"，分片："+prepend+" 删除失败", "\n        错误：", err)
	// 			} else {
	// 				fmt.Println("组件：" + component.Name + "，版本：" + version.Name + "，分片：" + prepend + " 删除成功")
	// 			}
	// 		}

	// 		os.Remove(fileName)
	// 	}

	// }

	url, err := utils.ComponentUrl(d.component, d.version)
	if err != nil {
		return
	}
	req := utils.NewRequest(url)

	remoteFile, err := req.Total()
	if err != nil {
		fmt.Println(d.component+"-"+d.version, "发送错误：", err)
		return
	}
	total := remoteFile.Total
	var process int = d.process
	if process == 0 {
		dm := total / 1024 / 1024
		process = utils.CalcProcess(dm)
	}
	var partitionTotal = total / int64(process)
	var fileName string
	if len(remoteFile.FileName) > 0 {
		fileName = remoteFile.FileName
	} else {
		urlObj, _ := netUrl.Parse(url)
		urlPath := urlObj.Path
		lastIdx := strings.LastIndex(urlPath, "/")
		fileName = string([]rune(urlPath)[lastIdx+1:])
	}
	for i := 0; i < process; i++ {
		start := partitionTotal * int64(i)
		end := partitionTotal*(int64(i)+1) - 1
		if process-1 == i {
			end = total
		}
		prepend := utils.Short(url + ":" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10))
		err = os.Remove(prepend)
		if err != nil {
			fmt.Println("组件："+d.component+"，分片："+prepend+" 删除失败", "\n        错误：", err)
		} else {
			fmt.Println("组件：" + d.component + "，分片：" + prepend + " 删除成功")
		}
	}

	os.Remove(fileName)

}
