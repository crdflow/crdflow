package template_funcs

import (
	"strings"
	"text/template"
)

var Func = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},

	"lower": func(s string) string {
		return strings.ToLower(s)
	},
}
