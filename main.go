package main

import (
	"downloader/config"
	"downloader/delete"
	"downloader/downloader"
	"downloader/show"
	"flag"
	"fmt"
	"os"
)

func main() {

	var isV bool
	var u string
	var c string
	var v string
	var p int
	var f string
	listFlagSet := flag.NewFlagSet("list", flag.ExitOnError)
	listFlagSet.BoolVar(&isV, "v", false, "查看版本")
	listFlagSet.StringVar(&c, "c", "", "查看组件")
	listFlagSet.StringVar(&f, "f", "", "指定外部配置文件")
	downloadFlagSet := flag.NewFlagSet("download", flag.ExitOnError)
	downloadFlagSet.StringVar(&u, "u", "", "URL下载")
	downloadFlagSet.StringVar(&c, "c", "", "组件下载")
	downloadFlagSet.StringVar(&v, "v", "", "具体组件版本下载")
	downloadFlagSet.IntVar(&p, "p", 0, "启用几个并发同步下载")
	downloadFlagSet.StringVar(&f, "f", "", "指定外部配置文件")
	deleteFlagSet := flag.NewFlagSet("list", flag.ExitOnError)
	deleteFlagSet.StringVar(&v, "v", "", "指定版本")
	deleteFlagSet.StringVar(&c, "c", "", "指定组件")
	deleteFlagSet.StringVar(&f, "f", "", "指定外部配置文件")
	deleteFlagSet.IntVar(&p, "p", 0, "启用几个并发同步下载")

	if len(os.Args) == 1 {
		config.Load(f)
		showList := show.NewShow(c, isV)
		showList.Show()
		return
	}

	switch os.Args[1] {
	case "list":
		listFlagSet.Parse(os.Args[2:])
		config.Load(f)
		showList := show.NewShow(c, isV)
		showList.Show()
	case "delete":
		deleteFlagSet.Parse(os.Args[2:])
		config.Load(f)
		delete := delete.NewDeleter(c, v, p)
		delete.Delete()
	case "download":
		downloadFlagSet.Parse(os.Args[2:])
		config.Load(f)
		var d *downloader.Downloader
		var err error
		if len(u) > 0 {
			d, err = downloader.NewUrl(u, p)
		} else {
			d, err = downloader.NewComponent(c, v, p)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		d.Download()
	default:
		listFlagSet.Parse(os.Args[1:])
		config.Load(f)
		showList := show.NewShow(c, isV)
		showList.Show()
	}

}
