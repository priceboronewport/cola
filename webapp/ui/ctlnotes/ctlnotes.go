package ctlnotes

import (
	"fmt"
	"github.com/priceboronewport/cola/elementw"
	"github.com/priceboronewport/cola/webapp/ui"
)

type CtlNotes struct {
	Label    *element.Element
	TextArea *element.Element
}

func New(pg *ui.Page, label string, id string, value string, rows int, cols int) *CtlNotes {
	n := CtlNotes{Label: element.New("label"), TextArea: element.New("textarea")}
	n.Label.Attributes["class"] = "ctl"
	n.Label.InnerHTML = label
	n.TextArea.Attributes["class"] = "ctl"
	n.TextArea.Attributes["id"] = id
	n.TextArea.Attributes["name"] = id
	n.TextArea.Attributes["rows"] = fmt.Sprintf("%d", rows)
	n.TextArea.Attributes["cols"] = fmt.Sprintf("%d", cols)
	n.TextArea.InnerHTML = value
	pg.AddStylesheet("/res/css/ctl.css")
	return &n
}

func (n *CtlNotes) OuterHTML() string {
	if n.Label.InnerHTML == "" {
		return n.TextArea.OuterHTML()
	}
	return n.Label.OuterHTML() + "<br/>" + n.TextArea.OuterHTML()
}
