package src_v2

const (
	TEMPLATES_DIR_RUST = `./templates_v2/`
	//OUTPUT_DIR_RUST    = `/hamid/life/_active/backbone/lib/shared/src/gen/my_play/`
	OUTPUT_DIR_RUST = `/hamid/life/_active/backbone/dev_play/tmp/src/my_dev/`
)

//var DATABASES = []string{"twitter"}

var DATABASES = []string{"flip_my"}

var triggerNeededArr = []string{}

var OutPutBuffer = &GenOut{}

var EscapeColumnNames = false
