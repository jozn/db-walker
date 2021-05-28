package src_v2

const (
	TEMPLATES_DIR_RUST = `./templates_v2`
	OUTPUT_DIR_RUST    = `/hamid/life/_active/backbone/lib/shared/src/gen/my_play/`
)

var DATABASES = []string{"twitter"}

var triggerNeededArr = []string{}

var OutPutBuffer = &GenOut{}

var EscapeColumnNames = false
