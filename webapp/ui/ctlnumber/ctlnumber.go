package ctlnumber

import (
	".."
	"../../../element"
)

type CtlText struct {
	Label *element.Element
	Input *element.Element
}

func New(pg *ui.Page, label string, id string, value string) *CtlText {
	ctl := CtlText{Label: element.New("label"), Input: element.New("input")}
	ctl.Label.Attributes["class"] = "ctl"
	ctl.Label.InnerHTML = label
	ctl.Input.Attributes["class"] = "ctl"
	ctl.Input.Attributes["style"] = "text-align: right;"
	ctl.Input.Attributes["maxlength"] = "14"
	ctl.Input.Attributes["size"] = "8"
	ctl.Input.Attributes["id"] = id
	ctl.Input.Attributes["name"] = id
	ctl.Input.Attributes["type"] = "text"
	pg.AddStylesheet("/res/css/ctl.css")
	pg.AddScript("/res/js/ctl.js")
	if value != "" {
		ctl.Input.Attributes["value"] = value
	}
	return &ctl
}

func (ctl *CtlText) OuterHTML() (html string) {
	if ctl.Label.InnerHTML != "" {
		html = ctl.Label.OuterHTML() + "<br/>"
	}
	html += ctl.Input.OuterHTML() + "<script type='text/javascript'>ctl.Filter(document.getElementById('" + ctl.Input.Attributes["id"] + "'), function(value) { return /^-?[0-9\\.\\,]*$/.test(value); });</script>"
	return
}
