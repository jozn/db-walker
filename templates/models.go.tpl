package {{.PackageName}}

{{range .Tables }}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .TableName }} '{{ .TableNameGo }}'.
{{- end }}
type {{ .TableNameGo }} struct {
{{- range .Columns }}
	{{ .ColumnName }} {{ .GoTypeOut }} {{ ms_col_comment_json .Comment }} {{ ms_col_comment_raw .Comment }}      {{/* `json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }} */}}
{{- end }}
{{- if .PrimaryKey }}
	{{/* // xox fields */}}
	_exists, _deleted bool
{{ end -}}
}
/*
:= &{{ .TableNameGo }} {
{{- range .Columns }}
	{{ .ColumnName }}: {{.GoDefaultOut}},
{{- end }}
*/
{{end}}
