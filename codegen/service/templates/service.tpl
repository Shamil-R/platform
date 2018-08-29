// Code generated, DO NOT EDIT.

package {{ .PackageName }}

import (
{{- range $import := .Imports }}
	{{- $import.Write }}
{{ end }}
)

type {{ .Name }}}Service interface {
    
}