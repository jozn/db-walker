package src

import (
	"bytes"
	"io/ioutil"
	"ms/sun/helper"
	"os"
	"os/exec"
	"text/template"
)

func build(gen *GenOut) {

	writeOutput("zz_xo.go", buildFromTemplate("xo.go.tpl", gen))
	writeOutput("models.go", buildFromTemplate("models.go.tpl", gen))
	genTablesOrma("orm.go.tpl", gen)

	/*OutGoRPCsStr := buildFromTemplate("rpc.tgo", gen)
	writeOutput("pb__gen_ant.go", OutGoRPCsStr)

	OutGoRPCsEmptyStr := buildFromTemplate("rpc_empty_imple.tgo", gen)
	writeOutput("pb__gen_ant_empty.go", OutGoRPCsEmptyStr)

	writeOutput("pb__gen_enum.proto", buildFromTemplate("enums.proto", gen))
	writeOutput("RPC_HANDLERS.java", buildFromTemplate("RPC_HANDLERS.java", gen))
	writeOutput("PBFlatTypes.java", buildFromTemplate("PBFlatTypes.java", gen))
	writeOutput("flat.go", buildFromTemplate("flat.tgo", gen))
	*/
	e1 := exec.Command("gofmt", "-w", OUTPUT_DIR).Run()
	e2 := exec.Command("goimports", "-w", OUTPUT_DIR).Run()
	helper.NoErr(e1)
	helper.NoErr(e2)
}

func genTablesOrma(tplName string, gen *GenOut) {
	tpl := _getTemplate(tplName)

	for _, table := range gen.Tables {
		buffer := bytes.NewBufferString("")
		err := tpl.Execute(buffer, table)
		helper.NoErr(err)
		writeOutput("z_"+table.TableName+".go", buffer.String())
	}

}

func writeOutput(fileName, output string) {
	println(output)
	ioutil.WriteFile(OUTPUT_DIR+fileName, []byte(output), os.ModeType)

}

func buildFromTemplate(tplName string, gen *GenOut) string {
	tpl := template.New("" + tplName)
	tpl.Funcs(NewTemplateFuncs())
	tplGoInterface, err := ioutil.ReadFile(TEMPLATES_DIR + tplName)
	helper.NoErr(err)
	tpl, err = tpl.Parse(string(tplGoInterface))
	helper.NoErr(err)

	buffer := bytes.NewBufferString("")
	err = tpl.Execute(buffer, gen)
	helper.NoErr(err)

	return buffer.String()
}

func _getTemplate(tplName string) *template.Template {
	tpl := template.New("" + tplName)
	tpl.Funcs(NewTemplateFuncs())
	tplGoInterface, err := ioutil.ReadFile(TEMPLATES_DIR + tplName)
	helper.NoErr(err)
	tpl, err = tpl.Parse(string(tplGoInterface))
	helper.NoErr(err)
	return tpl
}
