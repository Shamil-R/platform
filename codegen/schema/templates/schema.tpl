{{ range $s, $def := .Types }}
{{ template "item" $def }}
{{ end }}
