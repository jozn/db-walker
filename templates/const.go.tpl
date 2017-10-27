package xconst

const (
{{range .Tables}}
	{{.TableNameGo}}_Table = "{{.TableName}}"
	{{.TableNameGo}}_TableGo = "{{.TableNameGo}}"
{{- end}}
)
{{range .Tables}}
	var {{.TableNameGo}} = 	struct {
		{{range .Columns}}
			{{.ColumnName}} string
		{{- end}}
	}{
		{{range .Columns}}
        	{{.ColumnName}}: "{{.ColumnName}}",
        {{- end}}
	}
{{end}}

