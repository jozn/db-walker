package src

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
)

func rustBuild(gen *GenOut) {

	rustGenModels(gen)

	//rustWriteOutput("z_xo.go", rustBuildFromTemplate("xo.go.tpl", gen))

	if false {
		e1 := exec.Command("gofmt", "-w", OUTPUT_DIR_GO_X).Run()
		e2 := exec.Command("goimports", "-w", OUTPUT_DIR_GO_X).Run()
		NoErr(e1)
		NoErr(e2)
	}
}

func rustGenTablesOrma(tplName string, gen *GenOut) {
	tpl := _goGetTemplate(tplName)

	for _, table := range gen.Tables {
		buffer := bytes.NewBufferString("")
		err := tpl.Execute(buffer, table)
		NoErr(err)
		goWriteOutput("zz_"+table.TableName+".go", buffer.String())
	}
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
	tables := []*Table{}
	for _, t := range gen.Tables {
		if !skipTableModel(t.TableNameSql) {
			tables = append(tables, t)
		}
	}

	buffer := bytes.NewBufferString("")
	err := tpl.Execute(buffer, tables)
	NoErr(err)
	rustWriteOutput("models.rs", buffer.String())
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
