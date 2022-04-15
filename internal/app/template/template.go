package template

import (
	"html/template"
	"strings"
	"sync"
)

var (
	funcMap     []template.FuncMap
	funcMapOnce sync.Once
)

// FuncMap returns a list of user-defined template functions.
func FuncMap() []template.FuncMap {
	funcMapOnce.Do(func() {
		funcMap = []template.FuncMap{map[string]interface{}{
			"Safe": Safe,
		}}
	})
	return funcMap
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}

// NewLine2br simply replaces "\n" to "<br>".
func NewLine2br(raw string) string {
	return strings.Replace(raw, "\n", "<br>", -1)
}
