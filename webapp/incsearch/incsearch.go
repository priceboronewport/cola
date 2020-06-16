package incsearch

import (
	"../../element"
)

type IncSearch struct {
    Label *element.Element
	Input *element.Element
	List  *element.Element
}

func New(label string, id string, value string) *IncSearch {
	is := IncSearch{Label: element.New("label"), Input: element.New("input"), List: element.New("div")}
    is.Label.Attributes["for"] = id
    is.Label.InnerHTML = label
	is.Input.Attributes["onBlur"] = "is_blur(this)"
	is.Input.Attributes["onKeyup"] = "is_change(this)"
	is.Input.Attributes["type"] = "text"
	is.Input.Attributes["id"] = id
    is.Input.Attributes["name"] = id
    is.Input.Attributes["autocomplete"] = "off"
    is.Input.Attributes["value"] = value
	is.List.Attributes["id"] = id + "_list"
	is.List.Attributes["tabindex"] = "-1"
	is.List.Attributes["class"] = "list"
    is.List.Attributes["style"] = "height: 200px"
	return &is
}

func (is *IncSearch) OuterHTML() string {
	return is.Label.OuterHTML() + "<br/>" + is.Input.OuterHTML() + "<br/>" + is.List.OuterHTML()
}
