package config

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Components struct {
	XMLName    xml.Name    `xml:"Components"`
	Components []Component `xml:"Component"`
}
type Component struct {
	Name     string    `xml:"name,attr"`
	Base     string    `xml:"base,attr"`
	Versions []Version `xml:"Version"`
}
type Version struct {
	Name string `xml:"name,attr"`
	Url  string `xml:"url,attr"`
}

var (
	C = new(Components)
)

//go:embed data.xml
var dataConfig []byte

func Load(f string) {
	var buf []byte
	var err error
	if f != "" {
		buf, err = ioutil.ReadFile(f)
		if err != nil {
			buf = dataConfig
		}
	} else {
		buf = dataConfig
	}
	err = xml.Unmarshal(buf, C)
	if err != nil {
		fmt.Println(err)
	}
}
