package src_v2

import (
	"bytes"
	"io/ioutil"
	//"strings"
	"text/template"
)

func rustBuild_dep(gen *GenOut) {

	//rustGenModels_dep(gen)

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

func rustGenModels_dep(gen *GenOut) {
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

// todo (maybe): b/c of diffrence in api of cassandar and mysql libs for now we not support Ins > use or_{col} for now

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
