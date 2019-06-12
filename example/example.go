/*
 *  example - Sample web application using Cola webapp framework.
 *
 *  Copyright (c) 2019  Priceboro Newport, Inc.  All Rights Reserved.
 *
 */

package main

import (
	"../webapp"
	"html"
	"html/template"
	"net/http"
)

type RenderParams struct {
	Head template.HTML
	Body template.HTML
}

func main() {
	webapp.Register("", "/", Handler, true)
	webapp.ListenAndServe("./conf/")
}

func Handler(w http.ResponseWriter, r *http.Request, p webapp.HandlerParams) {
	head := "<title>Hello</title>"
	user_rec := webapp.User(p.Username)
	body := "<b><i>Hello " + html.EscapeString(user_rec["first_name"]) + "</i></b><hr/><a href='/logout'>Logout</a>"
	Render(w, head, body)
}

func Render(w http.ResponseWriter, head string, body string) {
	webapp.Render(w, "default.html", RenderParams{Head: template.HTML(head), Body: template.HTML(body)})
}
