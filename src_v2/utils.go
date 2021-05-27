package src_v2

import (
	"strings"
	"unicode"
)

func SnakeToCamel(s string) string {
	var r string

	if len(s) == 0 {
		return s
	}

	//ME: hack snake just for those of having "_"
	if strings.Index(s, "_") < 0 {
		return strings.ToUpper(s[:1]) + s[1:]
	}

	for _, w := range strings.Split(s, "_") {
		if w == "" {
			continue
		}

		r += strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
	}

	return r
}

// ToSnake convert the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

func sqlTypesToRustType(sqlType string) (typ, org, def string) {
	switch strings.ToLower(sqlType) {
	case "string", "varchar", "char", "text", "tinytext":
		typ = "String"
		org = "&str"
		def = `"".to_string()`
	case "bool", "boolean":
		typ = "bool"
		org = "bool"
		def = `false`
	case "tinyint", "smallint", "mediumint", "int", "integer":
		typ = "u32"
		org = "u32"
		def = `0u32`
	case "bigint":
		typ = "u64"
		org = "u64"
		def = `0u64`
	//case "bytes", "blob":
	//	typ = "Blob"
	//	org = "&Blob"
	//	def = `Blob::new(vec![])`
	case "binary", "blob", "mediumblob", "longblob":
		typ = "Vec<u8>"
		org = "&Vec<u8>"
		def = `vec![]`
	case "decimal":
		typ = "f64"
		org = "f64"
		def = `0f64`
	case "float":
		typ = "f32"
		org = "f32"
		def = `0f32`

	default:
		typ = "UNKNOWN_sqlToRust__" + sqlType
		org = "UNKNOWN_sqlToRust__" + sqlType
		def = `""`
	}
	//duration,timeuuid, uuid, map, tuple, set, list
	return
}
