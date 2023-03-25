package utils

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func ProxyIps() ([]string, error) {
	url := "https://proxy.seofangfa.com"
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		fmt.Printf("%v解析失败，原因：%v\n", url, err.Error())
		return nil, err
	}

	addresss := make([]string, 0)
	trNodes := doc.Find("tbody").Find("tr").Nodes
	for _, trNode := range trNodes {
		nextNode := trNode.FirstChild
		ip := nodetext(nextNode)
		nextNode = nextNode.NextSibling
		port := nodetext(nextNode)
		nextNode = nextNode.NextSibling
		nextNode = nextNode.NextSibling
		area := nodetext(nextNode)
		if strings.Index(area, "中国") == -1 {
			address := ip + ":" + port
			addresss = append(addresss, address)

		}

	}
	return addresss, nil
}

func ProxyIp() (string, error) {
	addresss, err := ProxyIps()
	if err != nil {
		return "", err
	}
	if len(addresss) == 0 {
		return "", errors.New("没有找到IP")
	}
	rand := GetRandomWithMin(0, len(addresss))
	address := addresss[rand]
	return address, nil
}

func nodetext(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	} else if node.FirstChild != nil {
		var buf bytes.Buffer
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			buf.WriteString(nodetext(c))
		}
		return buf.String()
	}

	return ""
}
