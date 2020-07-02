package ui

import (
	".."
	"html"
	"html/template"
	"net/http"
)

type Page struct {
	Head    string
	header  string
	Content string
}

type uiParams struct {
	Head    template.HTML
	Header  template.HTML
	Content template.HTML
	User    template.HTML
	Menu    template.HTML
}

func New(title, icon string) *Page {
	head := "<title>" + html.EscapeString(title) + "</title>"
	if icon != "" {
		head += webapp.Icon(icon)
		head += webapp.Stylesheet("/res/css/ui.css")
		head += webapp.Script("/res/js/ui.js")
	}
	pg := Page{Head: head}
	return &pg
}

func (pg *Page) AddHeaderIcon(src string, url string) {
	if url == "" {
		pg.header += "<td><div class='icon' style='background-image: url(\"" + src + "\")'></div></td>"
	}
}

func (pg *Page) AddHeaderLabel(label string, url string) {
	if url == "" {
		pg.header += "<td class='active'>" + html.EscapeString(label) + "</td>"
	} else {
		pg.header += "<td><a href='" + url + "'>" + html.EscapeString(label) + "</a></td>"
	}
}

func (pg *Page) Render(w http.ResponseWriter, r *http.Request, p webapp.HandlerParams) {
	header := "<table id='ui_header'><tr>"
	var menu string
	if p.Username != "" {
		header += "<td id='ui_menu'><button class='ui_menu_button' onClick='UIMenuShow(\"ui_menu_content\")'></button></td>"
		if webapp.HasPermission(p.Username, "switch_user") {
			menu += "<a href='/su'>Switch User</a><hr/>"
		}
		menu += "<a href='/logout'>Log Out</a>"
	}
        header += pg.header
	if p.Username != "" {
                user := webapp.User(p.Username)
		header += "<td id='ui_user'><button class='ui_menu_button' onClick='UIMenuShow(\"ui_usermenu_content\")'>" + user["first_name"] + " " + user["last_name"] + "</button></td>"
	}
	header += "</tr></table>"

	webapp.Render(w, "ui.html", uiParams{Head: template.HTML(pg.Head),
		Header: template.HTML(header), Content: template.HTML(pg.Content),
		User: template.HTML(p.Username), Menu: template.HTML(menu)})
}
