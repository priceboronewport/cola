package ctllist

import (
	".."
	"../../../element"
	"database/sql"
	"fmt"
)

type CtlList struct {
	Label  *element.Element
	Select *element.Element
}

func New(pg *ui.Page, label string, id string) *CtlList {
	cl := CtlList{Label: element.New("label"), Select: element.New("select")}
	cl.Label.Attributes["class"] = "ctl"
	cl.Label.InnerHTML = label
	cl.Select.Attributes["class"] = "ctl"
	cl.Select.Attributes["id"] = id
	cl.Select.Attributes["name"] = id
	pg.AddStylesheet("/res/css/ctl.css")
	return &cl
}

func (cl *CtlList) Add(label string, value string) {
	cl.Select.InnerHTML += "<option value='" + value + "'>" + label + "</option>"
}

func (cl *CtlList) Load(db *sql.DB, query string) (err error) {
	rows, err := db.Query(query)
	if err == nil {
		defer rows.Close()
		var label, value string
		for rows.Next() {
			rows.Scan(&label, &value)
			cl.Add(label, value)
		}
	} else {
		fmt.Printf(" ** ERROR: %s\n", err.Error())
	}
	return
}

func (cl *CtlList) OuterHTML() string {
	if cl.Label.InnerHTML == "" {
		return cl.Select.OuterHTML()
	}
	return cl.Label.OuterHTML() + "<br/>" + cl.Select.OuterHTML()
}

func (cl *CtlList) SetSelected(selected string) {
	if selected != "" {
		options, err := element.Parse(cl.Select.InnerHTML)
		if err == nil {
			html := ""
			for _, option := range options {
				if option.Tag == "option" && ((option.Attributes["value"] == selected) || (option.InnerHTML == selected)) {
					option.Attributes["selected"] = ""
					html += option.OuterHTML()
				} else {
					html += option.OuterHTML()
				}
			}
			cl.Select.InnerHTML = html
		}
	}
}
