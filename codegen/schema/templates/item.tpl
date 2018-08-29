{{- $createInput := cat .Name "CreateInput" | nospace -}}
{{- $updateInput := cat .Name "UpdateInput" | nospace -}}
{{- $whereUniqueInput := cat .Name "WhereUniqueInput" | nospace -}}
{{- $whereInput := cat .Name "WhereInput" | nospace }}

type {{ .Name }} {
    {{- range $, $field := .Fields }}
    {{ .Name }}: {{ .Type }}
    {{- end }}
}

input {{ $createInput }} {
    {{- range $, $field := .Fields }}
    {{- if and (ne .Type.Name "ID") (not .Type.Elem) }}
    {{ .Name }}: {{ .Type }}
    {{- end }}
    {{- end }}
}

input {{ $updateInput }} {
    {{- range $, $field := .Fields }}
    {{- if and (ne .Type.Name "ID") (not .Type.Elem) }}
    {{ .Name }}: {{ .Type }}
    {{- end }}
    {{- end }}
}

input {{ $whereUniqueInput }} {
    id: ID
}

input {{ $whereInput }} {
    id: ID
}

type Mutation {
    create{{ .Name }}(data: {{ $createInput }}!): {{ .Name }}!
    update{{ .Name }}(data: {{ $updateInput }}!, where: {{ $whereUniqueInput }}!): {{ .Name }}
    delete{{ .Name }}(where: {{ $whereUniqueInput }}!): {{ .Name }}
}

type Query {
    {{ lower .Name }}(where: {{ $whereUniqueInput }}!): {{ .Name }}
    {{ lower .Name | plural }}(where: {{ $whereInput }}): [{{ .Name }}]!
}
