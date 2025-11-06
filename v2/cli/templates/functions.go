package templates

import (
	"text/template"

	E "github.com/IBM/fp-go/v2/either"
)

var (
	templateFunctions = template.FuncMap{}
)

func Parse(name, tmpl string) E.Either[error, *template.Template] {
	return E.TryCatchError(template.New(name).Funcs(templateFunctions).Parse(tmpl))
}
