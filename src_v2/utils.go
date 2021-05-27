package src_v2

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/knq/snaker"
)

//copy of "ms/xox/snaker"

// SnakeToCamel converts s to CamelCase.
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
		/*u := strings.ToUpper(w)
		if ok := commonInitialisms[u]; ok {//me not need we use: Id and Html
			r += u
		} else {
			r += strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}*/
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

// commonInitialisms is the set of commonInitialisms.
//
// taken from: github.com/golang/lint @ 206c0f0
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

var strip = regexp.MustCompile(`\(.*\)`)

// This is newer version
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

var PrecScaleRE = regexp.MustCompile(`\(([0-9]+)(\s*,[0-9]+)?\)$`)

// SinguralizeIdentifier will singularize a identifier, returning it in
// CamelCase.
func SingularizeIdentifier(s string) string {
	/*if i := reverseIndexRune(s, '_'); i != -1 {
		s = s[:i] + "_" + inflector.Singularize(s[i+1:])
	} else {
		s = inflector.Singularize(s)
	}*/

	return snaker.SnakeToCamelIdentifier(s)
}
