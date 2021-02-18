package src

import (
	"fmt"
	"strings"
)

type GenOut struct {
	PackageName          string
	Tables               []*Table
	TablesTriggers       []*Table
	GeneratedPb          string
	GeneratedPbConverter string
}

type Table struct {
	TableName         string
	Columns           []*Column
	HasPrimaryKey     bool
	PrimaryKey        *Column
	DataBase          string //or schema in PG or tablesapce in cassandra
	Seq               int
	Comment           string
	IsAutoIncrement   bool
	Indexes           []*Index
	TableSchemeOut    string //with table "`ms`.`post`"
	TableNameSql      string //"post"
	TableNameGo       string
	TableNameJava     string
	TableNamePB       string
	ShortName         string
	NeedTrigger       bool   // MySql trigger events
	XPrimaryKeyGoType string //shortcut
	IsMysql           bool
	IsPG              bool   // Is PostgreSQL/CockroachDB
	Dollar            string // Use ? For MySql
}

func (t *Table) ColNum() int {
	return len(t.Columns)
}

type Column struct {
	ColumnName      string
	ColumnNameCamel string
	ColumnNameSnake string //not used
	SqlType         string // bigint(20) > NOT USED
	SqlTypeStrip    string // bigint
	Seq             int    // From 1
	Comment         string
	ColumnNameOut   string //dep: unclear what is the meaning
	GoTypeOut       string
	RoachTypeOut    string
	GoDefaultOut    string
	JavaTypeOut     string
	PBTypeOut       string
	StructTagOut    string // For Go Json tag
	IsNullAble      bool
	IsPrimary       bool
	IsUnique        bool
	//Rust
	RustTypeOut string
}

type Index struct {
	FuncNameOut string // index_name
	IndexName   string // index_name
	IsUnique    bool   // is_unique
	IsPrimary   bool   // is_primary
	SeqNo       int    // seq_no
	Columns     []*Column
	Table       *Table
}

//////////////////////////////////////

func (t *Table) GetColumnByName(col string) *Column {
	for _, c := range t.Columns {
		if c.ColumnName == col {
			return c
		}
	}
	return nil
}

func (t *Column) ToCockroachColumns() string {
	s := t.ColumnNameSnake + " " + t.RoachTypeOut + " "
	if t.IsPrimary {
		s += "PRIMARY KEY "
	}
	if t.IsUnique {
		s += "UNIQUE "
	}
	if !t.IsNullAble {
		s += "NOT NULL "
	}

	return s
}

// For Rust
func (t *Column) GetColIndex() int {
	return t.Seq - 1
}

func (t *Table) GetRustParam() string {
	arr := []string{}
	for _, c := range t.Columns {
		arr = append(arr, fmt.Sprintf("self.%s.clone().into()", c.ColumnName))
	}
	return strings.Join(arr, ", ")
}

func (t *Table) GetRustParamNoPrimaryKey() string {
	arr := []string{}
	for _, c := range t.Columns {
		if c.ColumnName != t.PrimaryKey.ColumnName {
			arr = append(arr, fmt.Sprintf("self.%s.clone().into()", c.ColumnName))
		}
	}
	return strings.Join(arr, ", ")
}
