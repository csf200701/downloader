package downloader

import (
	"downloader/bar"
	"downloader/config"
	"downloader/utils"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Downloader struct {
	url              string
	componentName    string
	componentVersion string
	process          int
	req              *utils.Request
}

type MergeFile struct {
	partition string
	start     int64
	end       int64
	status    int //1正常 2异常
}

func NewUrl(url string, process int) (*Downloader, error) {
	if len(url) == 0 {
		return nil, errors.New("URL地址不能为空")
	}
	req := utils.NewRequest(url)
	return &Downloader{
		url:     url,
		process: process,
		req:     req,
	}, nil
}

func NewComponent(componentName string, componentVersion string, process int) (*Downloader, error) {
	if len(componentName) == 0 {
		return nil, errors.New("组件不能为空")
	}
	var url string
	var base string
	c := config.C
	for _, component := range c.Components {
		if component.Name == componentName {
			if len(componentVersion) > 0 {
				for _, version := range component.Versions {
					if version.Name == componentVersion {
						url = version.Url
						base = component.Base
						break
					}
				}
			} else {
				url = component.Versions[0].Url
				base = component.Base
			}
			break
		}
	}
	if len(url) == 0 {
		return nil, errors.New("获取不到组件所对应的URL地址")
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
	req := utils.NewRequest(url)
	return &Downloader{
		url:              url,
		componentName:    componentName,
		componentVersion: componentVersion,
		process:          process,
		req:              req,
	}, nil
}

func (d *Downloader) Download() {
	total, err := d.req.Total()
	if err != nil {
		fmt.Println(err)
		return
	}
	dm := total / 1024 / 1024
	if d.process == 0 {
		if dm < 40 {
			d.process = 1
		} else if dm < 1024 {
			d.process = 5
		} else {
			d.process = 10
		}
	}
	var partitionTotal = total / int64(d.process)

	urlObj, _ := url.Parse(d.url)
	urlPath := urlObj.Path
	lastIdx := strings.LastIndex(urlPath, "/")
	fileName := string([]rune(urlPath)[lastIdx+1:])

	bc := bar.New(fmt.Sprintf("【%s】文件总大小：%v，分片数：%v，每个分片平均大小：%v", fileName, utils.CalcSize(total), d.process, utils.CalcSize(partitionTotal)))

	recieves := make([]chan *MergeFile, 0)
	for i := 0; i < d.process; i++ {
		start := partitionTotal * int64(i)
		end := partitionTotal*(int64(i)+1) - 1
		if d.process-1 == i {
			end = total
		}
		var c chan *MergeFile = make(chan *MergeFile, 1)
		recieves = append(recieves, c)
		go d.partitionDownload(c, bc, start, end)
	}

	partitionFiles := make([]*MergeFile, 0)
	for _, cm := range recieves {
		m := <-cm
		partitionFiles = append(partitionFiles, m)
	}

	mergeFile, _ := os.Create(fileName)
	defer mergeFile.Close()
	for _, m := range partitionFiles {
		buf, err := ioutil.ReadFile(m.partition)
		mergeFile.WriteAt(buf, m.start)
		err = os.Remove(m.partition)
		if err != nil {
			panic(1)
		}
	}

}

func (d *Downloader) partitionDownload(c chan *MergeFile, bc *bar.BarContainer, start, end int64) {
	size := end - start + 1
	prepend := utils.Short(d.url + ":" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10))
	b := bc.NewBar(prepend, size)
	fw := NewFileWriter(bc, b, prepend, start, end, size)
	defer fw.Close()

	var status int = 1
	if start+fw.offset < end {
		reader, realSize, _ := d.req.Content(start+fw.offset, end)
		defer reader.Close()
		b.Add((end-(start+fw.offset)+1)-realSize, 0)

		wLen, err := io.Copy(fw, reader)
		if err != nil || wLen != realSize {
			status = 2
		}
	}
	c <- &MergeFile{partition: prepend, start: start, end: end, status: status}
}
