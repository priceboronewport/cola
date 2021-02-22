package element

import (
	"errors"
	"fmt"
	"strings"
)

var void_tags = map[string]int{"img": 1, "input": 1, "hr": 1, "area": 1, "link": 1, "br": 1, "meta": 1, "base": 1, "col": 1, "embed": 1, "keygen": 1, "param": 1, "source": 1, "track": 1, "wbr": 1}

type Element struct {
	Tag        string
	Attributes map[string]string
	InnerHTML  string
}

func New(tag string) *Element {
	element := Element{Tag: strings.ToLower(tag), Attributes: make(map[string]string), InnerHTML: ""}
	return &element
}

func Parse(html string) (elements []Element, err error) {
	var str string
	var inner_html string
	var current_tag string
	in := false
	depth := 0
	for i, ch := range html {
		if ch == '<' {
			if !in {
				if len(str) > 0 {
					if depth == 0 { // text tag
						e := Element{InnerHTML: str}
						elements = append(elements, e)
					} else {
						inner_html += str
					}
				}
				in = true
				str = "<"
			} else {
				err = errors.New(fmt.Sprintf("Invalid '<' @ %d", i))
				return
			}
		} else if ch == '>' {
			if in {
				str += string(ch)
				if len(str) > 2 && str[0:2] == "</" {
					depth -= 1
					if depth == 0 {
						tag := ParseTag(current_tag)
						if tag == ParseTag(str) { // html tag
							a, er := ParseAttributes(current_tag[len(tag)+1 : len(current_tag)-1])
							if er != nil {
								err = er
								return
							}
							e := Element{Tag: ParseTag(current_tag), InnerHTML: inner_html, Attributes: a}
							inner_html = ""
							elements = append(elements, e)
						} else {
							err = errors.New(fmt.Sprintf("Mismatched end tag '%s' @ %d", str, i))
							return
						}
					} else {
						inner_html += str
					}
				} else {
					if depth == 0 {
						current_tag = str
						tag := ParseTag(current_tag)
						if void_tags[tag] == 1 { // void tag
							a, er := ParseAttributes(current_tag[len(tag)+1 : len(current_tag)-1])
							if er != nil {
								err = er
								return
							}
							e := Element{Tag: tag, Attributes: a}
							elements = append(elements, e)
							current_tag = ""
							depth -= 1
						} else if tag == "!--" { // comment tag
							e := Element{Tag: tag, InnerHTML: str[4 : len(str)-3]}
							elements = append(elements, e)
							current_tag = ""
							depth -= 1
						}
					} else {
						inner_html += str
					}
					depth += 1
				}
				str = ""
				in = false
			} else {
				err = errors.New(fmt.Sprintf("Invalid '>' @ %d", i))
				return
			}
		} else {
			str += string(ch)
		}
	}
	if len(str) > 0 {
		if str[0] == '<' {
			err = errors.New("Unclosed '<'")
		} else { // text tag
			e := Element{InnerHTML: str}
			elements = append(elements, e)
		}
	}
	return
}

func ParseAttributes(html string) (attributes map[string]string, err error) {
	attributes = make(map[string]string)
	state := 0
	var a, v string
	instr := ' '
	for _, ch := range html {
		if state == 0 {
			if ch != ' ' {
				a = string(ch)
				state = 1
			}
		} else if state == 1 {
			if ch == '=' {
				v = ""
				state = 2
			} else if ch == ' ' {
				attributes[a] = ""
				state = 0
			} else {
				a += string(ch)
			}
		} else if state == 2 {
			if v == "" {
				if ch == '\'' || ch == '"' {
					instr = ch
				} else {
					v += string(ch)
				}
			} else if ch == instr {
				attributes[a] = v
				a = ""
				v = ""
				instr = ' '
				state = 0
			} else {
				v += string(ch)
			}
		}
	}
	if state == 2 {
		attributes[a] = v
	}
	return
}

func ParseTag(html string) (tag string) {
	if len(html) > 4 && html[0:4] == "<!--" {
		return "!--"
	}
	in := false
	for _, ch := range strings.ToLower(html) {
		if ch >= 'a' && ch <= 'z' {
			if !in {
				in = true
			}
			tag += string(ch)
		} else {
			if in {
				return
			}
		}
	}
	return
}

func (element *Element) OuterHTML() (html string) {
	if element.Tag != "" {
		html = "<" + element.Tag
		if element.Tag == "!--" {
			html += element.InnerHTML + "--"
		} else {
			if len(element.Attributes) > 0 {
				for k, v := range element.Attributes {
					html += " " + k
					if v != "" {
						html += "='" + v + "'"
					}
				}
			}
			if void_tags[element.Tag] == 1 {
				html += "/"
			} else {
				html += ">" + element.InnerHTML + "</" + element.Tag
			}
		}
		html += ">"
	} else {
		html = element.InnerHTML
	}
	return
}
