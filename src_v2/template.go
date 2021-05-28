package src_v2

import (
	"fmt"
	"strings"
)
import "text/template"

func NewTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"toUpper":        strings.ToUpper,
		"colnames":       colnames,
		"colvals_dollar": colvals_dollar,
	}
}

// colnames creates a list of the column names found in fields, excluding any
// Field with Name contained in ignoreNames.
//
// Used to present a comma separated list of column names, that can be used in
// a SELECT, or UPDATE, or other SQL clause requiring an list of identifiers
// (ie, "field_1, field_2, field_3, ...").
func colnames(columns []*OutColumn, ignoreNames ...string) string {
	ignore := map[string]bool{}
	for _, n := range ignoreNames {
		ignore[n] = true
	}

	str := ""
	i := 0
	for _, f := range columns {
		if ignore[f.ColumnName] {
			continue
		}

		if i != 0 {
			str = str + ", "
		}
		str = str + colname(f)
		i++
	}

	return str
}

// colname returns the ColumnName of col, optionally escaping it if
// ArgType.EscapeColumnNames is toggled.
func colname(col *OutColumn) string {
	if EscapeColumnNames {
		//return c.Loader.Escape(ColumnEsc, col.ColumnName)
		return fmt.Sprintf("`%s`", col.ColumnName)
	}

	return col.ColumnName
}

func colvals_dollar(table *OutTable, fields []*OutColumn, ignoreNames ...string) string {
	ignore := map[string]bool{}
	for _, n := range ignoreNames {
		ignore[n] = true
	}

	str := ""
	i := 0
	for _, f := range fields {
		if ignore[f.ColumnName] {
			continue
		}

		if i != 0 {
			str = str + ", "
		}
		dollar := "?"
		/*		if table.IsPG {
				dollar = fmt.Sprintf("$%d", i+1)
			}*/
		str = str + dollar //c.Loader.NthParam(i)
		i++
	}

	return str
}
