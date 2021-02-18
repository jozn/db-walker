package src

const (
	TEMPLATES_DIR_GO      = `./templates/go/`
	OUTPUT_DIR_GO_X       = `./out/shared/x/`
	OUTPUT_DIR_GO_X_CONST = `./out/shared/x/xconst/`
	OUTPUT_PROTO_DIR      = `./out/shared/proto/`

	TEMPLATES_DIR_RUST = `./templates/rust/`
	OUTPUT_DIR_RUST    = `/hamid/life/_active/backbone/micro/tmp/src/`
	//OUTPUT_DIR_RUST       = `./out/rust/shared/x/`

	FORMAT = true
)

var DATABASES = []string{"twitter"}

var DATABASES_COCKROACHE = []string{"suncdb"}

var triggerNeededArr = []string{}

var OutPutBuffer = &GenOut{
	PackageName: "x",
}

var EscapeColumnNames = false

// Old

//var DATABASES = []string{"sun", "sun_chat", "sun_file","sun_meta", "sun_push", "sun_log","sun_internal"}
//var triggerNeededArr = []string{"action","user", "chat", "post", "comment", "tags"}
