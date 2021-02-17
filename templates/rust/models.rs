use mysql_async::prelude::*;
use mysql_async::{FromRowError, OptsBuilder, Params, Row};
use mysql_common::row::ColumnIndex;

{{range . }}

#[derive(Debug, PartialEq, Eq, Clone)]
struct {{ .TableNameJava }}  { // {{ .TableName }}
{{- range .Columns }}
    {{ .ColumnName }}: {{ .RustTypeOut }},
{{- end }}
}

impl FromRow for {{ .TableNameJava }} {
    fn from_row_opt(row: Row) -> Result<Self, FromRowError>
    where
        Self: Sized,
    {
        Ok({{ .TableNameJava }}  {
        {{- range .Columns }}
            {{ .ColumnName }}: row.get({{ .GetColIndex }}).unwrap(),
        {{- end }}
        })
    }
}
{{end}}