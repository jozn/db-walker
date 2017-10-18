package src

const (
	TEMPLATES_DIR    = `C:\Go\_gopath\src\ms\snake\templates\`
	//OUTPUT_DIR       = `C:\Go\_gopath\src\ms\snake\play\`
	OUTPUT_DIR       = `C:\Go\_gopath\src\ms\sun\models\x\`
	OUTPUT_PROTO_DIR = `C:\Go\_gopath\src\ms\sun\models\protos\`
	DATABASE = "ms"
)

var OutPutBuffer = &GenOut{
	PackageName: "x",
}

var EscapeColumnNames = false
