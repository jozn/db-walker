{{- range $key,$table := .Tables }}
{{- with $table }}
/*Table: {{ .TableName }}  */
CREATE TABLE IF NOT EXISTS {{.TableSchemeOut}} (
	{{- range .Columns }}
    {{ .ToCockroachColumns }},
	{{- end}}
);
{{end -}}
{{end}}
}