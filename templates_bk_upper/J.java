package com.mardomsara.social.json;

public class J {

{{- range $key,$model := .Tables }}
{{- with $model }}
	public static class {{.TableNameJava}} {//oridnal: {{$key }}
		{{- range .Columns }}
		public {{ .JavaTypeOut }} {{ .ColumnName }};
		{{- end }}
	}
{{end -}}
{{end}}
}