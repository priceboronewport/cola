package ctlfile

import (
	"github.com/priceboronewport/cola/element"
	"github.com/priceboronewport/cola/webapp/ui"
)

type CtlFile struct {
	Label *element.Element
	File  *element.Element
}

func New(pg *ui.Page, label string, id string) *CtlFile {
	ctl := CtlFile{Label: element.New("label"), File: element.New("input")}
	ctl.Label.Attributes["class"] = "ctl"
	ctl.Label.InnerHTML = label
	ctl.File.Attributes["class"] = "ctl"
	ctl.File.Attributes["id"] = id
	ctl.File.Attributes["name"] = id
	ctl.File.Attributes["type"] = "file"
	pg.AddStylesheet("/res/css/ctl.css")
	return &ctl
}

func (ctl *CtlFile) OuterHTML() string {
	if ctl.Label.InnerHTML == "" {
		return ctl.File.OuterHTML()
	}
	return ctl.Label.OuterHTML() + "<br/>" + ctl.File.OuterHTML()
}
