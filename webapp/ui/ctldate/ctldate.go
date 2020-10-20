package ctldate

import (
	".."
	"../../../element"
)

type CtlDate struct {
	Label *element.Element
	Input *element.Element
}

func New(pg *ui.Page, label string, id string, value string) *CtlDate {
	ctl := CtlDate{Label: element.New("label"), Input: element.New("input")}
	ctl.Label.Attributes["class"] = "ctl"
	ctl.Label.InnerHTML = label
	ctl.Input.Attributes["class"] = "ctl tcal"
    if id != "" {
	  ctl.Input.Attributes["id"] = id
	  ctl.Input.Attributes["name"] = id
    }
    if value != "" {
      ctl.Input.Attributes["value"] = value
    }
	ctl.Input.Attributes["type"] = "text"
	ctl.Input.Attributes["size"] = "10"
	ctl.Input.Attributes["maxlength"] = "10"
	pg.AddStylesheet("/res/css/ctl.css")
	pg.AddStylesheet("/res/css/ctldate.css")
	pg.AddScript("/res/js/ctldate.js")
	return &ctl
}

func (ctl *CtlDate) OuterHTML() string {
	if ctl.Label.InnerHTML == "" {
		return ctl.Input.OuterHTML()
	}
	return ctl.Label.OuterHTML() + "<br/>" + ctl.Input.OuterHTML()
}
