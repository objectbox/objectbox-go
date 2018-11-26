package templates

import (
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"StringTitle": strings.Title,
	"StringCamel": func(s string) string {
		result := strings.Title(s)
		return strings.ToLower(result[0:1]) + result[1:]
	},
}
