package main

import (
	"downloader/config"
	"downloader/deleter"
	"downloader/downloader"
	"downloader/show"
	"flag"
	"fmt"
	"os"
)

func main() {
	// for i := 0; i < 20; i++ {
	// 	address, _ := utils.ProxyIp()
	// 	fmt.Println(address)
	// }
	// return

	var isV bool
	var u string
	var c string
	var v string
	var p int
	var f string
	var px bool
	var pxh string
	var pxu string
	var pxp string
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
	downloadFlagSet.BoolVar(&px, "px", false, "是否使用代理请求")
	downloadFlagSet.StringVar(&pxh, "pxh", "", "使用代理请求的Host，格式：xxx.xxx.xxx.xxx:8080")
	downloadFlagSet.StringVar(&pxu, "pxu", "", "使用代理请求的用户名")
	downloadFlagSet.StringVar(&pxp, "pxp", "", "使用代理请求的用户密码")
	deleteFlagSet := flag.NewFlagSet("list", flag.ExitOnError)
	deleteFlagSet.StringVar(&u, "u", "", "指定URL下载")
	deleteFlagSet.StringVar(&v, "v", "", "指定版本")
	deleteFlagSet.StringVar(&c, "c", "", "指定组件")
	deleteFlagSet.StringVar(&f, "f", "", "指定外部配置文件")
	deleteFlagSet.IntVar(&p, "p", 0, "启用几个并发同步下载")
	deleteFlagSet.BoolVar(&px, "px", false, "是否使用代理请求")
	deleteFlagSet.StringVar(&pxh, "pxh", "", "使用代理请求，格式：xxx.xxx.xxx.xxx:8080")
	deleteFlagSet.StringVar(&pxu, "pxu", "", "使用代理请求的用户名")
	deleteFlagSet.StringVar(&pxp, "pxp", "", "使用代理请求的用户密码")

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
		var d *deleter.Deleter
		var err error
		if len(u) > 0 {
			d, err = deleter.NewUrlWithProxy(u, p, pxh, pxu, pxp, px)
		} else {
			d, err = deleter.NewComponentWithProxy(c, v, p, pxh, pxu, pxp, px)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		d.Delete()
	case "download":
		downloadFlagSet.Parse(os.Args[2:])
		config.Load(f)
		var d *downloader.Downloader
		var err error
		if len(u) > 0 {
			d, err = downloader.NewUrlWithProxy(u, p, pxh, pxu, pxp, px)
		} else {
			d, err = downloader.NewComponentWithProxy(c, v, p, pxh, pxu, pxp, px)
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
