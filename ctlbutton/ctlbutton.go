package ctlbutton

import (
	"../element"
)

func New(hint string, class string, url string) *element.Element {
	ctl := element.New("button")
	ctl.Attributes["class"] = class
	ctl.Attributes["onClick"] = "window.location.href=\"" + url + "\""
	if hint != "" {
		ctl.InnerHTML = "<span class='ui_hinttext'>" + hint + "</span>"
	}
	return ctl
}
