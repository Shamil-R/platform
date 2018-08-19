# types
{{- range $_, $def := .Types }}
type {{ .Name }} {
    {{- range $_, $field := $def.Fields }}
    {{ $field.Name }}: {{ $field.Describe }}
    {{- end }}
}
{{ end }}

{{- range $_, $def := .Types }}
# {{ $def.Name }} inputs

{{- with $def.Input.Create }}
input {{ .Name }} {
    {{- range $_, $field := .Fields }}
    {{ $field.Name }}: {{ $field.Describe }}
    {{- end }}
}
{{ end }}

{{- with $def.Input.Update }}
input {{ .Name }} {
    {{- range $_, $field := .Fields }}
    {{ $field.Name }}: {{ $field.Describe }}
    {{- end }}
}
{{ end }}

{{- with $def.Input.WhereUnique }}
input {{ .Name }} {
    id: ID!
}
{{ end }}

{{- with $def.Input.Where }}
input {{ .Name }} {
    id: ID!
}
{{ end }}

{{- end }}

# mutations
type Mutation {
{{- range $_, $def := .Types }}
{{- with $def }}
    # {{ .Name }} mutations
    {{ .Mutation.Create }}(data: {{ .Input.Create.Name }}): {{ .Name }}
    {{ .Mutation.Update }}(data: {{ .Input.Update.Name }}, where: {{ .Input.WhereUnique.Name }}): {{ $def.Name }}
    {{ .Mutation.Delete }}(where: {{ .Input.WhereUnique.Name }}): {{ .Name }}
{{- end }}
{{- end }}
}

# queries
type Query {
{{- range $_, $def := .Types }}
{{- with $def }}
    # {{ .Name }} queries
    {{ .Query.Item }}(where: {{ .Input.WhereUnique.Name }}!): {{ .Name }}
    {{ .Query.List }}(where: {{ .Input.Where.Name }}): [{{ .Name }}]!
{{- end }}
{{- end }}
}