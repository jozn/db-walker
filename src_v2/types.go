package src_v2

import (
	"fmt"
	"strings"
)

type GenOut struct {
	Tables         []*OutTable
	TablesFiltered []*OutTable // Only those table we have interest on
}

// Mysql Native Types

type NativeTable struct {
	TableName             string
	Columns               []*NativeColumn
	HasPrimaryKey         bool
	IsCompositePrimaryKey bool
	SinglePrimaryKey      *NativeColumn
	PrimaryKeys           []*NativeColumn //used for composite keys -- Note: not used in gen as
	DataBase              string
	Comment               string
	IsAutoIncrement       bool
	Indexes               []*NativeIndex
}

type NativeColumn struct {
	ColumnName      string
	SqlType         string // bigint varchar
	SqlTypeFull     string // bigint(20) or varchar(75) > NOT USED
	Ordinal         int    // From 1
	Comment         string
	IsNullAble      bool
	IsPrimary       bool // Multi columns could be primary
	IsUnique        bool // This shows if the columns has an unique index (applicable just to single column) -- NOT primary column in here
	IsAutoIncrement bool // A none primary column could be auto increment
}

type NativeIndex struct {
	IndexName string // index_name
	IsUnique  bool   // is_unique
	IsPrimary bool   // is_primary
	ColNum    int    // is_primary
	Comment   string
	Columns   []*NativeColumn
	Table     *NativeTable
}

// Output Types

type OutTable struct {
	TableName           string
	HasPrimaryKey       bool
	HasMultiPrimaryKeys bool
	IsAutoIncr          bool
	IsAutoIncrPrimary   bool
	DataBase            string
	Comment             string
	Columns             []*OutColumn
	SinglePrimaryKey    *OutColumn // deprecated
	AutoIncrKey         *OutColumn
	PrimaryKeys         []*OutColumn //used for composite keys -- Note: not used in gen as
	Indexes             []*OutIndex
	// Views
	SchemeTable    string // with database "`ms`.`post`"
	TableNameCamel string
}

type OutColumn struct {
	ColumnName      string
	Ordinal         int // From 1
	IsNullAble      bool
	IsSinglePrimary bool
	IsInPrimary     bool // deprecated if multi primary and this col is included
	IsPrimary       bool // if unique index or sinlge primary
	IsUnique        bool // if unique index or sinlge primary
	IsAutoIncr      bool
	// Views
	RustType       string
	RustTypeBorrow string

	WhereModifiersRust    []WhereModifier
	WhereInsModifiersRust []WhereModifierIns
}

// Only None Single IsPrimary
type OutIndex struct {
	IndexName string // index_name
	IsUnique  bool   // is_unique
	IsPrimary bool   // Just multi columns primary
	ColNum    int    // is_primary
	Columns   []*OutColumn
}

// For Rust
func (t *OutColumn) GetColIndex() int {
	return t.Ordinal - 1
}

func (t *OutTable) ColNum() int {
	return len(t.Columns)
}

func (t *NativeTable) GetColumnByName(col string) *NativeColumn {
	for _, c := range t.Columns {
		if c.ColumnName == col {
			return c
		}
	}
	return nil
}

func (t *OutTable) GetColumnByName(col string) *OutColumn {
	for _, c := range t.Columns {
		if c.ColumnName == col {
			return c
		}
	}
	return nil
}

func (t *OutColumn) IsNumber() bool {
	res := false
	switch t.RustType {
	case "u64", "i64", "u32", "i32", "f32", "f64":
		res = true
	}
	return res
}

func (t *OutColumn) IsString() bool {
	res := false
	switch t.RustType {
	case "String":
		res = true
	}
	return res
}

func (t *OutTable) GetRustParam() string {
	arr := []string{}
	for _, c := range t.Columns {
		arr = append(arr, fmt.Sprintf("self.%s.clone().into()", c.ColumnName))
	}
	return strings.Join(arr, ", ")
}

func (t *OutTable) GetRustParamNoPrimaryKey_dep() string {
	arr := []string{}
	for _, c := range t.Columns {
		if t.SinglePrimaryKey != nil && c.ColumnName != t.SinglePrimaryKey.ColumnName {
			arr = append(arr, fmt.Sprintf("self.%s.clone().into()", c.ColumnName))
		}
	}
	return strings.Join(arr, ", ")
}

func (t *OutTable) GetRustParamNoneIncrKeys() string {
	arr := []string{}
	for _, c := range t.Columns {
		if !c.IsAutoIncr {
			arr = append(arr, fmt.Sprintf("self.%s.clone().into()", c.ColumnName))
		}
	}
	return strings.Join(arr, ", ")
}

func (t *OutTable) GetRustParamPrimaryKeys() string {
	arr := []string{}
	for _, c := range t.Columns {
		if c.IsPrimary {
			arr = append(arr, fmt.Sprintf("self.%s.clone().into()", c.ColumnName))
		}
	}
	return strings.Join(arr, ", ")
}

func (t *OutTable) GetRustParamNoPrimaryKeys() string {
	arr := []string{}
	for _, c := range t.Columns {
		if !c.IsPrimary {
			arr = append(arr, fmt.Sprintf("self.%s.clone().into()", c.ColumnName))
		}
	}
	return strings.Join(arr, ", ")
}

func (t *OutTable) GetRustUpdateFrag() string {
	arr := []string{}
	for _, c := range t.Columns {
		if !c.IsPrimary {
			arr = append(arr, fmt.Sprintf("%s = ?", c.ColumnName))
		}
		/*if c.ColumnName != t.SinglePrimaryKey.ColumnName {
			arr = append(arr, fmt.Sprintf("%s = ?", c.ColumnName))
		}*/
	}
	return strings.Join(arr, ", ")
}

func (t *OutTable) GetRustUpdateKeysWhereFrag() string {
	arr := []string{}
	for _, c := range t.PrimaryKeys {
		arr = append(arr, fmt.Sprintf("%s = ?", c.ColumnName))
	}
	return strings.Join(arr, " AND ")
}

func (c *OutColumn) GetModifiersRust() (res []WhereModifier) {
	add := func(m WhereModifier) {
		if len(m.AndOr) > 0 {
			m.FuncName = m.Prefix + "_" + c.ColumnName + m.Suffix
		} else {
			m.FuncName = c.ColumnName + m.Suffix
		}
		res = append(res, m)
	}

	for _, andOr := range []string{"", "AND", "OR"} {
		if c.IsNumber() || c.IsString() {
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

func (c *OutColumn) GetRustModifiersIns() (res []WhereModifierIns) {
	add := func(m WhereModifierIns) {
		if len(m.AndOr) > 0 {
			m.FuncName = m.Prefix + "_" + c.ColumnName + m.Suffix
		} else {
			m.FuncName = c.ColumnName + m.Suffix
		}
		res = append(res, m)
	}

	for _, andOr := range []string{"", "AND", "OR"} {
		if c.IsNumber() || c.IsString() {
			add(WhereModifierIns{"_in", strings.ToLower(andOr), andOr, ""})
		}
	}

	return
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
