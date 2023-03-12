package show

import (
	"downloader/config"
	"fmt"
)

type ShowList struct {
	component string
	isVersion bool
}

func NewShow(component string, isVersion bool) *ShowList {
	return &ShowList{component: component, isVersion: isVersion}
}

func (s *ShowList) Show() {
	c := config.C
	var components []config.Component = make([]config.Component, 0)
	for _, component := range c.Components {
		isShow := false
		if len(s.component) > 0 {
			if s.component == component.Name {
				isShow = true
			}
		} else {
			isShow = true
		}
		if isShow {
			if !s.isVersion {
				component.Versions = make([]config.Version, 0)
			}
			components = append(components, component)
		}
	}
	//打印
	showContent := ""
	for i, component := range components {
		isHasVersion := false
		showContent += "  " + component.Name
		if len(component.Versions) > 0 {
			showContent += "：\n"
			isHasVersion = true
			for j, version := range component.Versions {
				showContent += "    " + version.Name
				if j < len(component.Versions)-1 {
					showContent += "\n"
				}
			}
		}
		if i < len(components)-1 {
			if isHasVersion {
				showContent += "\n"
			}
			showContent += "\n"
		}
	}
	if len(showContent) > 0 {
		showContent = "查询可用的下载软件如下：\n" + showContent
	}
	fmt.Println(showContent)
}
