package src

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
)

func goBuild(gen *GenOut) {

	goGenModels(gen)
	goWriteOutput("z_xo.go", goBuildFromTemplate("xo.go.tpl", gen))
	//goWriteOutput("z_models.go", goBuildFromTemplate("models.go.tpl", gen))
	goWriteOutput("z_cache.go", goBuildFromTemplate("cache.go.tpl", gen))
	goWriteOutput("z_event.go", goBuildFromTemplate("event.go.tpl", gen))
	goWriteOutput("z_manual.go", goBuildFromTemplate("manual.go", gen))
	goWriteOutput("z_index.go", goBuildFromTemplate("index.go.tpl", gen))
	goWriteOutput("z_cache_secondary_index.go", goBuildFromTemplate("cache_secondary_index.go.tpl", gen))
	goWriteOutput("J.java", goBuildFromTemplate("J.java", gen))
	goWriteOutput("triggers.sql", goBuildFromTemplate("triggers.sql", gen))
	goWriteOutput("trigger.go", goBuildFromTemplate("trigger.go.tpl", gen))

	goWriteOutput("_tables_lowers.sql", goBuildFromTemplate("_tables_lowers.sql", gen))
	goWriteOutput("_tables_to_cockroach.sql", goBuildFromTemplate("_tables_to_cockroach.sql", gen))

	goWriteOutputConst("tables.go", goBuildFromTemplate("const.go.tpl", gen))

	genTablesOrma("orm.go.tpl", gen)

	PtMsgdef, converter := Gen_ProtosForTables(gen.Tables)
	goWriteOutput("TablePBCon.go", converter)
	ioutil.WriteFile(OUTPUT_PROTO_DIR+"pb_tables.proto", []byte(PtMsgdef), os.ModeType)

	if false && FORMAT {
		e1 := exec.Command("gofmt", "-w", OUTPUT_DIR_GO_X).Run()
		e2 := exec.Command("goimports", "-w", OUTPUT_DIR_GO_X).Run()
		NoErr(e1)
		NoErr(e2)
	}
}

func genTablesOrma(tplName string, gen *GenOut) {
	tpl := _getTemplate(tplName)

	for _, table := range gen.Tables {
		buffer := bytes.NewBufferString("")
		err := tpl.Execute(buffer, table)
		NoErr(err)
		goWriteOutput("zz_"+table.TableName+".go", buffer.String())
	}

}

func goWriteOutput(fileName, output string) {
	//println(output)
	os.MkdirAll(OUTPUT_DIR_GO_X, 0777)
	ioutil.WriteFile(OUTPUT_DIR_GO_X+fileName, []byte(output), os.ModeType)

}

func goWriteOutputConst(fileName, output string) {
	//println(output)
	os.MkdirAll(OUTPUT_DIR_GO_X_CONST, 0777)
	ioutil.WriteFile(OUTPUT_DIR_GO_X_CONST+fileName, []byte(output), os.ModeType)

}

func goBuildFromTemplate(tplName string, gen *GenOut) string {
	tpl := template.New("" + tplName)
	tpl.Funcs(NewTemplateFuncs())
	tplGoInterface, err := ioutil.ReadFile(TEMPLATES_DIR_GO + tplName)
	NoErr(err)
	tpl, err = tpl.Parse(string(tplGoInterface))
	NoErr(err)

	buffer := bytes.NewBufferString("")
	err = tpl.Execute(buffer, gen)
	NoErr(err)

	return buffer.String()
}

func goGenModels(gen *GenOut) {
	tpl := _getTemplate("models.go.tpl")
	tables := []*Table{}
	for _, t := range gen.Tables {
		if !skipTableModel(t.TableNameSql) {
			tables = append(tables, t)
		}
	}

	buffer := bytes.NewBufferString("")
	err := tpl.Execute(buffer, tables)
	NoErr(err)
	goWriteOutput("z_models.go", buffer.String())
}

func _getTemplate(tplName string) *template.Template {
	tpl := template.New("" + tplName)
	tpl.Funcs(NewTemplateFuncs())
	tplGoInterface, err := ioutil.ReadFile(TEMPLATES_DIR_GO + tplName)
	NoErr(err)
	tpl, err = tpl.Parse(string(tplGoInterface))
	NoErr(err)
	return tpl
}
