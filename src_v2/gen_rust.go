package src_v2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

func rustBuild(gen *GenOut) {

	rustGenWithDatabase(gen, "mod.rs", "mod.rs")

	rustGenEachModel(gen)
}

func rustGenEachModel(gen *GenOut) {
	tpl := _rustGetTemplate("model.rs")
	type tableOutput struct {
		table     string
		genOutput string
	}

	var outs = []*tableOutput{}
	for _, table := range gen.TablesFiltered {
		buffer := bytes.NewBufferString("")
		err := tpl.Execute(buffer, table)
		NoErr(err)

		o := &tableOutput{
			table:     table.TableName,
			genOutput: buffer.String(),
		}
		outs = append(outs, o)
	}

	for _, o := range outs {
		rustWriteOutput(fmt.Sprintf("%s.rs", o.table), o.genOutput)
	}

}

func rustGenWithDatabase(gen *GenOut, templateFile string, outputFile string) {
	tpl := _rustGetTemplate(templateFile)
	buffer := bytes.NewBufferString("")
	err := tpl.Execute(buffer, gen.TablesFiltered)
	NoErr(err)
	rustWriteOutput(outputFile, buffer.String())
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

func rustWriteOutput(fileName, output string) {
	os.MkdirAll(OUTPUT_DIR_RUST, 0777)
	ioutil.WriteFile(OUTPUT_DIR_RUST+fileName, []byte(output), os.ModeType)
}
