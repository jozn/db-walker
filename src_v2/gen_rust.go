package src_v2

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func rustBuild(gen *GenOut) {

	rustGenModels(gen)

}

func rustWriteOutput(fileName, output string) {
	os.MkdirAll(OUTPUT_DIR_RUST, 0777)
	ioutil.WriteFile(OUTPUT_DIR_RUST+fileName, []byte(output), os.ModeType)
}

func rustBuildFromTemplate(tplName string, gen *GenOut) string {
	tpl := template.New("" + tplName)
	tpl.Funcs(NewTemplateFuncs())
	tplGoInterface, err := ioutil.ReadFile(TEMPLATES_DIR_RUST + tplName)
	NoErr(err)
	tpl, err = tpl.Parse(string(tplGoInterface))
	NoErr(err)

	buffer := bytes.NewBufferString("")
	err = tpl.Execute(buffer, gen)
	NoErr(err)

	return buffer.String()
}

func rustGenModels(gen *GenOut) {
	tpl := _rustGetTemplate("models.rs")
	tables := []*OutTable{}
	for _, t := range gen.Tables {
		// We can skip any tables that we do not want in here. For now process all of them.
		// todo support multi primay keys
		if t.SinglePrimaryKey == nil {
			continue
		}
		tables = append(tables, t)
	}

	buffer := bytes.NewBufferString("")
	err := tpl.Execute(buffer, tables)
	NoErr(err)
	rustWriteOutput("mysql_models.rs", buffer.String())
}

func _rustGetTemplate(tplName string) *template.Template {
	tpl := template.New("" + tplName)
	tpl.Funcs(NewTemplateFuncs())
	tplGoInterface, err := ioutil.ReadFile(TEMPLATES_DIR_RUST + tplName)
	NoErr(err)
	tpl, err = tpl.Parse(string(tplGoInterface))
	NoErr(err)
	return tpl
}

///////////////////////// Migration From cassandra-walker ///////////////

func (table *OutTable) GetRustWheresTmplOut() string {
	const TPL = `
    pub fn {{ .Mod.FuncName }} (&mut self, val: {{ .Col.RustTypeBorrow }} ) -> &mut Self {
        let w = WhereClause{
            condition: "{{ .Mod.AndOr }} {{ .Col.ColumnName }} {{ .Mod.Condition }} ?".to_string(),
            args: val.into(),
        };
        self.wheres.push(w);
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

		for j := 0; j < len(col.WhereModifiersRust); j++ {
			wmr := col.WhereModifiersRust[j]

			parm := struct {
				Table *OutTable
				Mod   WhereModifier
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

// todo (maybe): b/c of diffrence in api of cassandar and mysql libs for now we not support Ins > use or_{col} for now
func (table *OutTable) GetRustWhereInsTmplOut() string {
	const TPL = `
    pub fn {{ .Mod.FuncName }} (&mut self, val: Vec<{{ .Col.RustTypeBorrow }}> ) -> &mut Self {
		let len = val.len();
        if len == 0 {
            return self
        }

        let mut marks = "?,".repeat(len);
        marks.remove(marks.len()-1);
        let w = WhereClause{
			condition: format!("{{ .Mod.AndOr }} {{ .Col.ColumnName }} IN ({})", marks),
            args: val.into(),
        };
        self.wheres.push(w);
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
		self.order_by.push("{{ .Col.ColumnName }} ASC");
        self
    }

	pub fn order_by_{{ .Col.ColumnName }}_desc(&mut self) -> &mut Self {
		self.order_by.push("{{ .Col.ColumnName }} DESC");
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

			fnStr := rawTemplateOutput(TPL, parm)
			//fmt.Println(fnStr)
			fnsOut = append(fnsOut, fnStr)
		}
	}

	return strings.Join(fnsOut, "")
}

func rawTemplateOutput(templ string, data interface{}) string {
	tpl := template.New("fns")
	tpl, err := tpl.Parse(templ)
	NoErr(err)

	buffer := bytes.NewBufferString("")
	err = tpl.Execute(buffer, data)
	NoErr(err)
	outPut := buffer.String()
	return outPut
}
