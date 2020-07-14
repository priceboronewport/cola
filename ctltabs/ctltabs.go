package ctltabs

import (
	"../element"
	"../webapp/ui"
)

type CtlTabs struct {
	id      string
	url     string
	Active  string
	Tabs    *element.Element
	Content *element.Element
}

func New(pg *ui.Page, id string, url string) *CtlTabs {
	ctl := CtlTabs{id: id, url: url, Tabs: element.New("table"), Content: element.New("div")}
	ctl.Content.Attributes["id"] = id + "_content"
	pg.AddScript("/res/js/stdlib.js")
	pg.AddScript("/res/js/ctltabs.js")
	pg.AddStylesheet("/res/css/ctltabs.css")
	return &ctl
}

func (ctl *CtlTabs) AddTab(id string, label string) {
	tr, err := element.Parse(ctl.Tabs.InnerHTML)
	if err == nil {
		var tds string
		if len(tr) > 0 && tr[0].Tag == "tr" {
			tds = tr[0].InnerHTML
		}
		tds += "<td id='" + id + "' onClick='ctltabs.LoadContent(this,\"" + ctl.url + "\")'>" + label + "</td>"
		ctl.Tabs.InnerHTML = "<tr>" + tds + "</tr>"
	}
}

func (ctl *CtlTabs) OuterHTML() (html string) {
	html = "<div id='" + ctl.id + "' class='tabs_container'>" + ctl.Tabs.OuterHTML() + ctl.Content.OuterHTML() + "</div>"
	if ctl.Active != "" {
		html += "<script type='text/javascript'>ctltabs.LoadContent(document.getElementById('" + ctl.Active + "'),'" + ctl.url + "');</script>"
	}
	return
}
