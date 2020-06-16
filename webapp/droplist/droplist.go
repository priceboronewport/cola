package droplist

import (
	"../../element"
	"database/sql"
)

type DropList struct {
	Label  *element.Element
	Select *element.Element
}

func New(label string, id string) *DropList {
	dl := DropList{Label: element.New("label"), Select: element.New("select")}
	dl.Label.InnerHTML = label
	dl.Select.Attributes["id"] = id
	dl.Select.Attributes["name"] = id
	return &dl
}

func (dl *DropList) Add(label string, value string) {
	dl.Select.InnerHTML += "<option value='" + value + "'>" + label + "</option>"
}

func (dl *DropList) Load(db *sql.DB, query string) (err error) {
	rows, err := db.Query(query)
	if err == nil {
		defer rows.Close()
		var label, value string
		for rows.Next() {
			rows.Scan(&label, &value)
			dl.Add(label, value)
		}
	}
	return
}

func (dl *DropList) OuterHTML() string {
	return dl.Label.OuterHTML() + "<br/>" + dl.Select.OuterHTML()
}

func (dl *DropList) SetSelected(selected string) {
	if selected != "" {
		options, err := element.Parse(dl.Select.InnerHTML)
		if err == nil {
			html := ""
			for _, option := range options {
				if option.Tag == "option" && option.Attributes["value"] == selected {
					option.Attributes["selected"] = ""
					html += option.OuterHTML()
				} else {
					html += option.OuterHTML()
				}
			}
			dl.Select.InnerHTML = html
		}
	}
}
