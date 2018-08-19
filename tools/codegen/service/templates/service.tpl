package {{ .PackageName }}

import (
{{- range $import := .Imports }}
	{{- $import.Write }}
{{ end }}
)

{{ range $interface := .Interfaces }}
    {{ with $interface.Name }}
        type {{ . }} interface {
            Create{{ . }}() error
            Update{{ . }}() error
            Delete{{ . }}() error
        }
    {{ end }}
{{ end }}