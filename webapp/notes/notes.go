package notes

import (
	"../../element"
	"fmt"
)

type Notes struct {
	Label    *element.Element
	TextArea *element.Element
}

func New(label string, id string, value string, rows int, cols int) *Notes {
	n := Notes{Label: element.New("label"), TextArea: element.New("textarea")}
	n.Label.InnerHTML = label
	n.TextArea.Attributes["id"] = id
	n.TextArea.Attributes["name"] = id
	n.TextArea.Attributes["rows"] = fmt.Sprintf("%d", rows)
	n.TextArea.Attributes["cols"] = fmt.Sprintf("%d", cols)
	n.TextArea.InnerHTML = value
	return &n
}

func (n *Notes) OuterHTML() string {
	return n.Label.OuterHTML() + "<br/>" + n.TextArea.OuterHTML()
}
