package main

import (
	"downloader/downloader"
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
	listFlagSet := flag.NewFlagSet("list", flag.ExitOnError)
	listFlagSet.BoolVar(&isV, "v", false, "查看版本")
	downloadFlagSet := flag.NewFlagSet("download", flag.ExitOnError)
	downloadFlagSet.StringVar(&u, "u", "", "URL下载")
	downloadFlagSet.StringVar(&c, "c", "", "组件下载")
	downloadFlagSet.StringVar(&v, "v", "", "具体组件版本下载")
	downloadFlagSet.IntVar(&p, "p", 0, "启用几个并发同步下载")

	if len(os.Args) < 1 {
		fmt.Println("期望的命令list或download没有选择")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listFlagSet.Parse(os.Args[2:])
	case "download":
		downloadFlagSet.Parse(os.Args[2:])
		downloader := downloader.NewUrl(u, p)
		downloader.Download()
	default:
		listFlagSet.Parse(os.Args[1:])
	}

}

// func main() {
// 	// bc := bar.New("ddddd")
// 	// total := 1024 * 1024 * 20
// 	// step := 1024 * 1024
// 	// b1 := bc.NewBar("sss1", total)
// 	// b2 := bc.NewBar("sss2", total)
// 	// for i := 0; i < total; i += step {
// 	// 	go b1.Add(step)
// 	// 	go b2.Add(step)
// 	// 	time.Sleep(100 * 3 * time.Millisecond)
// 	// }
// 	fmt.Println(utils.Md5("MYSQL>0-100000"))
// 	fmt.Println(utils.Md4("MYSQL>0-100000"))
// 	fmt.Println(utils.Sha256("MYSQL>0-100000"))
// 	fmt.Println(utils.Sha512("MYSQL>0-100000"))
// 	fmt.Println(utils.Short("MYSQL>0-100000"))
// 	fmt.Println(utils.Short("MYSQLsdsdsdsdsdsdsver45gh>0-100000"))
// 	fmt.Println(utils.Short("哈希包对此很有帮助。请注意，这是对特定哈希实现的抽象。在软件包子目录中可以找到一些现成的。iiiiiiiiiiiiiiiiiiiii"))
// 	fmt.Println(utils.Short("哈希包对此很有帮助。请注意，这是对特中可以找到一些现成的。iiiiiiiiiiiiiiiiiiiii"))
// 	fmt.Println(utils.Short("要闻福州13℃习近平在参加江苏代表团审议时强调 牢牢把握高质量发展这个首要任务| 十四届全国人大一次会议在京开幕"))
// 	fmt.Println(utils.Short("政府工作报告的背后：“安欣”和冬奥冠军都曾受邀为报告提"))
// 	fmt.Println(utils.Short("胡和平已任中央宣传部分管日常工作的副部长"))
// 	fmt.Println(utils.Short("GDP增长5%左右，今年发展主要预期目标有这些"))
// 	fmt.Println(utils.Short("《我想和你唱》开播首期就玩大的！王心凌苏有朋联合造就青春回忆杀"))
// 	fmt.Println(utils.Short("同样是短视频爆红，把赵雷和李荣浩一对比，差距就出来了"))
// }
