package src

import (
	"bytes"
	"io/ioutil"
	"ms/sun/helper"
	"os"
	"text/template"
)

const (
	TEMPLATES_DIR = `C:\Go\_gopath\src\ms\snake\templates\`
	OUTPUT_DIR    = `C:\Go\_gopath\src\ms\snake\play\`
)

func build(gen *GenOut) {
	OutGoRPCsStr := buildFromTemplate("rpc.tgo", gen)
	writeOutput("pb__gen_ant.go", OutGoRPCsStr)

	OutGoRPCsEmptyStr := buildFromTemplate("rpc_empty_imple.tgo", gen)
	writeOutput("pb__gen_ant_empty.go", OutGoRPCsEmptyStr)

	writeOutput("pb__gen_enum.proto", buildFromTemplate("enums.proto", gen))
	writeOutput("RPC_HANDLERS.java", buildFromTemplate("RPC_HANDLERS.java", gen))
	writeOutput("PBFlatTypes.java", buildFromTemplate("PBFlatTypes.java", gen))
	writeOutput("flat.go", buildFromTemplate("flat.tgo", gen))

}

func buildFromTemplate(tplName string, gen *GenOut) string {
	tpl := template.New("go_interface" + tplName)
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

func writeOutput(fileName, output string) {
	ioutil.WriteFile(OUTPUT_DIR+fileName, []byte(output), os.ModeType)
}
