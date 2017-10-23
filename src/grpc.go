package src

import (
	"bytes"
	"strings"
	"text/template"
)

type ProtoFile struct {
	Messages []ProtoMessageDef
	FileName string
	OutPut   string
}

type ProtoMessageDef struct {
	Fields           []ProtoMessageFieldDef
	MessageName      string
	IsTableNotInline bool //for complex types like user we need to refrence fields with dot (.)
}

type ProtoMessageFieldDef struct {
	TagId   int
	Name    string
	TypeMix SqlToPBType
	PBType  string
	Repeat  bool
}

func Gen_ProtosForTables(tbls []*Table) (protoMsgDefs, protoConv string) {
	filePB := ProtoFile{FileName: "pb_tables"}

	for _, t := range tbls {
		tpb := ProtoMessageDef{
			MessageName:      t.TableNamePB,
			IsTableNotInline: skipTableModel(t.TableName),
		}

		for i, f := range t.Columns {
			fpb := ProtoMessageFieldDef{
				TagId:   (i*2 + 1),
				Name:    f.ColumnName,
				TypeMix: MysqlParseTypeToProtoclBuffer(f.SqlType, true),
				Repeat:  false,
			}
			tpb.Fields = append(tpb.Fields, fpb)
		}

		filePB.Messages = append(filePB.Messages, tpb)
	}

	//gen proto def
	tmpl, err := template.New("t").Parse(TMP_PB)
	if err != nil {
		panic(err)
	}
	out := bytes.NewBufferString("")
	err = tmpl.Execute(out, filePB)
	if err != nil {
		panic(err)
	}
	OutPutBuffer.GeneratedPb += out.String()

	//gen proto conv
	tmplCon, err := template.New("t").Parse(TMP_PB_CONVERTER)
	if err != nil {
		panic(err)
	}
	outConv := bytes.NewBufferString("")
	err = tmplCon.Execute(outConv, filePB)
	if err != nil {
		panic(err)
	}
	OutPutBuffer.GeneratedPbConverter = outConv.String()
	//return outConv.String()
	//fmt.Println("size of PB (tabels) : ", len(c.Loader.CacheTables))
	return OutPutBuffer.GeneratedPb, OutPutBuffer.GeneratedPbConverter
}

var GRPC_TYOPES_MAP = map[string]string{ // go type to => PB types
	"int":     "int64",
	"string":  "string",
	"float32": "float",
	"float64": "double",
}

const TMP_PB = `
syntax = "proto3";
option java_package = "ir.ms.pb";
option java_multiple_files = true;
option optimize_for = LITE_RUNTIME; //CODE_SIZE;

option go_package = "x";

{{range .Messages}}
message PB_{{.MessageName }} {
    {{- range .Fields }}
    {{.TypeMix.PB }} {{.Name}} = {{.TagId}};
    {{- end }}
}

{{- end}}
`

const TMP_PB_CONVERTER = `
package x

{{range .Messages}}
/*
func PBConvPB__{{.MessageName }}_To_{{.MessageName }}( o *PB_{{.MessageName }}) *{{.MessageName }} {
  {{- if .IsTableNotInline -}}
   n := &{{.MessageName}}{}
    {{- range .Fields }}
   n.{{.Name}} = {{.TypeMix.Go}} ( o.{{.Name}} )
    {{- end -}}

  {{else }}
     n := &{{.MessageName}}{
    {{- range .Fields }}
      {{.Name}}: {{.TypeMix.Go}} ( o.{{.Name}} ),
      {{- end }}
    }
  {{- end }}
    return n
}

func PBConvPB_{{.MessageName }}_To_{{.MessageName }} ( o *{{.MessageName }}) *PB_{{.MessageName }} {
  {{- if .IsTableNotInline -}}
   n := &PB_{{.MessageName}}{}
    {{- range .Fields }}
   n.{{.Name}} = {{.TypeMix.GoGen}} ( o.{{.Name}} )
    {{- end -}}

  {{else }}
     n := &PB_{{.MessageName}}{
    {{- range .Fields }}
      {{.Name}}: {{.TypeMix.GoGen}} ( o.{{.Name}} ),
      {{- end }}
    }
  {{- end }}
    return n
}
*/
{{- end}}
`

//============================================================

type SqlToPBType struct {
	Go    string //simple go
	GoGen string //go type from pb genrator
	table string
	Java  string
	PB    string
}

//cp of MyParseType(...)
func MysqlParseTypeToProtoclBuffer(dt string, fromMysql bool) SqlToPBType {
	precision := 0
	unsigned := false

	res := SqlToPBType{}
	// extract unsigned

	if fromMysql {
		if strings.HasSuffix(dt, " unsigned") {
			unsigned = true
			dt = dt[:len(dt)-len(" unsigned")]
		}

		// extract precision
		dt, precision, _ = ParsePrecision(dt)
		_ = precision
		_ = unsigned
	}

	switch dt {
	case "bool", "boolean":
		res = SqlToPBType{
			Go:    "bool",
			GoGen: "bool",
			table: "bool",
			Java:  "Boolean",
			PB:    "bool",
		}

	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		res = SqlToPBType{
			Go:    "string",
			GoGen: "string",
			table: "text",
			Java:  "String",
			PB:    "string",
		}

	case "tinyint", "smallint", "mediumint", "int", "integer":
		res = SqlToPBType{
			Go:    "int",
			GoGen: "int32",
			table: "int",
			Java:  "Integer",
			PB:    "int32",
		}

	case "bigint":
		//the main diffrence is for int64
		res = SqlToPBType{
			Go:    "int",
			GoGen: "int64",
			table: "bigint",
			Java:  "Long",
			PB:    "int64",
		}

	case "float":
		res = SqlToPBType{
			Go:    "float32",
			GoGen: "float32",
			table: "float",
			Java:  "Float",
			PB:    "float",
		}

	case "decimal", "double":
		res = SqlToPBType{
			Go:    "float64",
			GoGen: "float64",
			table: "double",
			Java:  "Double",
			PB:    "double",
		}

	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		res = SqlToPBType{
			Go:    "[]byte",
			GoGen: "[]byte",
			table: "binary",
			Java:  "[]byte",
			PB:    "bytes",
		}

	case "timestamp", "datetime", "date", "time":

	default:
	}
	return res
}
