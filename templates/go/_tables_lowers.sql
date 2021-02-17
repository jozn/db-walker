{{- range $key,$table := .Tables }}
{{- with $table }}
/*Table: {{ .TableName }}  */
	{{- range .Columns }}
ALTER TABLE {{ $table.TableSchemeOut }} CHANGE COLUMN {{ .ColumnName }} {{ .ColumnNameSnake }} {{ .SqlType }};
	{{- end}}
{{end -}}
{{end}}
}