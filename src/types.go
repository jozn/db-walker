package src

import (
	"fmt"
	"strings"
)

type GenOut struct {
	PackageName          string
	Tables               []*Table
	RustTables               []*Table
	TablesTriggers       []*Table
	GeneratedPb          string
	GeneratedPbConverter string
}

// Note due to how our policy change from Go ro Rust, there could be some wired behaviour around primary keys
type Table struct {
	TableName         string
	Columns           []*Column
	HasPrimaryKey     bool
	IsCompositePrimaryKey     bool
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
	// Rust
	TableNameRust string
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
	ColumnNameRust string
	RustTypeOut string
	TypeRustBorrow string //todo
	WhereModifiersRust []WhereModifier //todo
	WhereInsModifiersRust []WhereModifierIns //todo

}

type Index struct {
	FuncNameOut string // index_name
	IndexName   string // index_name
	IsUnique    bool   // is_unique
	IsPrimary   bool   // is_primary
	SeqNo       int    // seq_no
	Columns     []*Column
	Table       *Table
	// Rust
	RustFuncName string // index_name
}

type WhereModifier struct {
	Suffix    string
	Prefix    string
	Condition string
	AndOr     string
	FuncName  string
}

type WhereModifierIns struct {
	Suffix   string
	Prefix   string
	AndOr    string
	FuncName string
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

func (t *Column) IsNumber() bool {
	res := false
	switch t.RustTypeOut {
	case "u64", "i64","u32", "i32", "f32", "f64":
		res = true
	}
	return res
}

func (t *Column) IsString() bool {
	res := false
	switch t.GoTypeOut {
	case "string":
		res = true
	}
	return res
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

func (t *Table) GetRustUpdateFrag() string {
	arr := []string{}
	for _, c := range t.Columns {
		if c.ColumnName != t.PrimaryKey.ColumnName {
			arr = append(arr, fmt.Sprintf("%s = ?", c.ColumnName))
		}
	}
	return strings.Join(arr, ", ")
}


////////////// Modifer for Rust /////////////

func (c *Column) GetModifiersRust() (res []WhereModifier) {
	add := func(m WhereModifier) {
		if len(m.AndOr) > 0 {
			m.FuncName = m.Prefix + "_" + c.ColumnNameRust + m.Suffix
		} else {
			m.FuncName = c.ColumnNameRust + m.Suffix
		}
		res = append(res, m)
	}

	for _, andOr := range []string{"", "AND", "OR"} {
		if c.IsNumber() || c.IsString(){
			pre := strings.ToLower(andOr)

			add(WhereModifier{"_eq", strings.ToLower(andOr), "=", andOr, ""})

			add(WhereModifier{"_lt", pre, "<", andOr, ""})
			add(WhereModifier{"_le", pre, "<=", andOr, ""})
			add(WhereModifier{"_gt", pre, ">", andOr, ""})
			add(WhereModifier{"_ge", pre, ">=", andOr, ""})
		}
	}

	return
}

func (c *Column) GetRustModifiersIns() (res []WhereModifierIns) {
	add := func(m WhereModifierIns) {
		if len(m.AndOr) > 0 {
			m.FuncName = m.Prefix + "_" + c.ColumnNameRust + m.Suffix
		} else {
			m.FuncName = c.ColumnNameRust + m.Suffix
		}
		res = append(res, m)
	}

	for _, andOr := range []string{"", "AND", "OR"} {
		if c.IsNumber() || c.IsString(){
			add(WhereModifierIns{"_in", strings.ToLower(andOr), andOr, ""})
		}
	}

	return
}
