package ctlbutton

import (
	"github.com/priceboronewport/cola/webapp/ui"
	"github.com/priceboronewport/cola/element"
)

func New(pg *ui.Page, hint string, class string, url string) *element.Element {
	ctl := element.New("button")
	ctl.Attributes["class"] = class
	ctl.Attributes["onClick"] = "window.location.href=\"" + url + "\""
	if hint != "" {
		ctl.InnerHTML = "<span class='ui_hinttext'>" + hint + "</span>"
	}
	pg.AddStylesheet("/res/css/ctl.css")
	return ctl
}
