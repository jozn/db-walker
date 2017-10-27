package src

const (
	TEMPLATES_DIR = `C:\Go\_gopath\src\ms\snake\templates\`
	//OUTPUT_DIR_GO_X       = `C:\Go\_gopath\src\ms\snake\play\`
	OUTPUT_DIR_GO_X  = `C:\Go\_gopath\src\ms\sun\models\x\`
	OUTPUT_DIR_GO_X_CONST  = `C:\Go\_gopath\src\ms\sun\models\x\xconst\`
	OUTPUT_PROTO_DIR = `C:\Go\_gopath\src\ms\sun\models\protos\`
	DATABASE         = "ms"

	FORMAT = true
)

var OutPutBuffer = &GenOut{
	PackageName: "x",
}

var EscapeColumnNames = false
