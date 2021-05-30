package src_v2

import (
	"bytes"
	"strings"
	"text/template"
)

///////////////////////// Migration From cassandra-walker ///////////////

func (table *OutTable) GetRustWheresTmplOut() string {
	const TPL = `
    pub fn {{ .Mod.FuncName }} (&mut self, val: {{ .Col.RustTypeBorrow }} ) -> &mut Self {
        let w = WhereClause{
            condition: "{{ .Mod.AndOr }} {{ .Col.ColumnName }} {{ .Mod.Condition }} ?".to_string(),
            args: val.into(),
        };
        self.q.wheres.push(w);
        self
    }
`

	fnsOut := []string{}

	for i := 0; i < len(table.Columns); i++ {
		col := table.Columns[i]

		for j := 0; j < len(col.WhereModifiersRust); j++ {
			wmr := col.WhereModifiersRust[j]

			parm := struct {
				Table *OutTable
				Mod   WhereModifier
				Col   *OutColumn
			}{
				table, wmr, col,
			}

			fnStr := _runPartialTmpl("GetRustWheresTmplOut_tmpl", TPL, parm)
			//fmt.Println(fnStr)
			fnsOut = append(fnsOut, fnStr)
		}
	}

	return strings.Join(fnsOut, "")
}

func (table *OutTable) GetRustWhereInsTmplOut() string {
	const TPL = `
    pub fn {{ .Mod.FuncName }} (&mut self, val: Vec<{{ .Col.RustTypeBorrow }}> ) -> &mut Self {
		let len = val.len();
        if len == 0 {
            return self
        }

        let mut marks = "?,".repeat(len);
        marks.remove(marks.len()-1);

		let arr = val.iter().map(|v|v.into()).collect();

        let w = WhereInClause{
			condition: format!("{{ .Mod.AndOr }} {{ .Col.ColumnName }} IN ({})", marks),
            args: arr,
        };
        self.q.wheres_ins.push(w);
        self
    }
`
	fnsOut := []string{}

	// parse template
	tpl := template.New("fns")
	tpl, err := tpl.Parse(TPL)
	NoErr(err)

	for i := 0; i < len(table.Columns); i++ {
		col := table.Columns[i]

		for j := 0; j < len(col.WhereInsModifiersRust); j++ {
			wmr := col.WhereInsModifiersRust[j]

			parm := struct {
				Table *OutTable
				Mod   WhereModifierIns
				Col   *OutColumn
			}{
				table, wmr, col,
			}

			buffer := bytes.NewBufferString("")
			err = tpl.Execute(buffer, parm)

			fnStr := buffer.String()
			//fmt.Println(fnStr)
			fnsOut = append(fnsOut, fnStr)
		}
	}

	return strings.Join(fnsOut, "")
}

// Selectors
func (table *OutTable) GetRustSelectorOrders() string {
	const TPL = `
    pub fn order_by_{{ .Col.ColumnName }}_asc(&mut self) -> &mut Self {
		self.q.order_by.push("{{ .Col.ColumnName }} ASC");
        self
    }

	pub fn order_by_{{ .Col.ColumnName }}_desc(&mut self) -> &mut Self {
		self.q.order_by.push("{{ .Col.ColumnName }} DESC");
        self
    }
`
	fnsOut := []string{}

	for i := 0; i < len(table.Columns); i++ {
		col := table.Columns[i]
		if col.IsNumber() || col.IsString() { //&& col.IsNumber()
			parm := struct {
				Table *OutTable
				Col   *OutColumn
			}{
				table, col,
			}

			fnStr := _runPartialTmpl("GetRustSelectorOrders", TPL, parm)
			//fmt.Println(fnStr)
			fnsOut = append(fnsOut, fnStr)
		}
	}

	return strings.Join(fnsOut, "")
}

func _runPartialTmpl(tplName string, templ string, data interface{}) string {
	tpl := template.New(tplName)
	tpl, err := tpl.Parse(templ)
	NoErr(err)

	buffer := bytes.NewBufferString("")
	err = tpl.Execute(buffer, data)
	NoErr(err)
	outPut := buffer.String()
	return outPut
}
