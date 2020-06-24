package ctlfile

import (
	"../element"
)

type CtlFile struct {
	Label *element.Element
	File  *element.Element
}

func New(label string, id string) *CtlFile {
	ctl := CtlFile{Label: element.New("label"), File: element.New("input")}
	ctl.Label.InnerHTML = label
	ctl.File.Attributes["id"] = id
	ctl.File.Attributes["name"] = id
	ctl.File.Attributes["type"] = "file"
	return &ctl
}

func (ctl *CtlFile) OuterHTML() string {
	return ctl.Label.OuterHTML() + "<br/>" + ctl.File.OuterHTML()
}
