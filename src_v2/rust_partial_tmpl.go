package src_v2

import (
	"bytes"
	"fmt"
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
func (table *OutTable) GetRustOrdersTmplOut() string {
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



// This produces Getter by primary key [whether this functionality was worth the effort is a good question]
func (oIndex *OutIndex) GetRustPrimaryGetter(table *OutTable) string {
	// Primary Getters
	type _RustGetter_ struct {
		//fnName    string
		paramName string
		paramType string
		callName  string
	}

	//============== Fill Array ===============
	arr := []_RustGetter_{}
	paramCnt := 1
	// Partion keys
	for i := 0; i < len(oIndex.Columns); i++ {
		col := oIndex.Columns[i]

		param := fmt.Sprintf("%s", col.ColumnName)
		//fnName := col.ColumnNameRust
		callName := fmt.Sprintf("%s_eq(%s)", col.ColumnName, param)
		if i > 0 {
			//fnName = fmt.Sprintf("_and_%s", col.ColumnNameRust)
			callName = fmt.Sprintf("and_%s_eq(%s)", col.ColumnName, param)
		}

		f := _RustGetter_{
			//fnName:    fnName,
			paramName: param,
			paramType: col.RustTypeBorrow,
			callName:  callName,
		}
		arr = append(arr, f)
		paramCnt += 1
	}

	//================ Make Str Output =================
	fnName := oIndex.IndexName
	if oIndex.IsPrimary {
		fnName = fmt.Sprintf("get_%s", table.TableName)
	}
	fnParam := []string{}
	fnSetter := fmt.Sprintf("%sSelector::new()", table.TableNameCamel)
	for i := 0; i < len(arr); i++ {
		f := arr[i]
		fnParam = append(fnParam, f.paramName+": "+f.paramType)
		fnSetter += "\n\t\t." + f.callName
	}

	//================ Build ==================
	const TPL_UNIQUE = `
pub async fn {{ .FnName }}({{ .FnParam }}, spool: &SPool) -> Result<{{.Table.TableNameCamel}},MyError> {
	let m = {{ .FnSetter }}
		.get_row(spool).await?;
	Ok(m)
}
`

	const TPL_MULTI = `
pub async fn {{ .FnName }}({{ .FnParam }}, spool: &SPool) -> Result<Vec<{{.Table.TableNameCamel}}>,MyError> {
	let m = {{ .FnSetter }}
		.get_rows(spool).await?;
	Ok(m)
}
`
	parm := struct {
		FnName   string
		FnParam  string
		FnSetter string
		Table    *OutTable
	}{
		FnName:   fnName,
		FnParam:  strings.Join(fnParam, ","),
		FnSetter: fnSetter,
		Table:    table,
	}

	if oIndex.IsUnique {
		out := _runPartialTmpl("INDX", TPL_UNIQUE, parm)
		return out
	}

	out := _runPartialTmpl("INDX", TPL_MULTI, parm)
	return out
}
/*
// Updater
func (table *OutTable) GetRustUpdaterFnsOut() string {
	const TPL = `
    pub fn update_{{ .Col.ColumnNameRust }}(&mut self, val: {{ .Col.TypeRustBorrow }}) -> &mut Self {
        self.updates.insert("{{ .Col.ColumnName }} = ?", val.into());
        self
    }
`

	const TPL_BLOB = `
    pub fn update_{{ .Col.ColumnNameRust }}(&mut self, val: {{ .Col.TypeRustBorrow }}) -> &mut Self {
        self.updates.insert("{{ .Col.ColumnName }} = ?", Blob::new(val.clone()).into());
        self
    }
`
	fnsOut := []string{}

	for i := 0; i < len(table.Columns); i++ {
		col := table.Columns[i]

		parm := struct {
			Table *OutTable
			Col   *OutColumn
		}{
			table, col,
		}

		var fnStr string

		// Due to cdrs lib limitation we should treat blob differently
		if col.TypeCql == "blob" {
			fnStr = rawTemplateOutput(TPL_BLOB, parm)
		} else {
			fnStr = rawTemplateOutput(TPL, parm)
		}

		//fmt.Println(fnStr)
		fnsOut = append(fnsOut, fnStr)
	}

	return strings.Join(fnsOut, "")
}
*/
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
