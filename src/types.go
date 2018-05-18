package src

type GenOut struct {
	PackageName          string
	Tables               []*Table
	TablesTriggers       []*Table
	GeneratedPb          string
	GeneratedPbConverter string
}

//dep
type DataBase struct {
	Tables []Table
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
	NeedTrigger       bool
	XPrimaryKeyGoType string //shortcut
	IsMysql           bool
	IsPG              bool
	Dollar            string
}

func (t *Table) ColNum() int {
	return len(t.Columns)
}

type Column struct {
	ColumnName      string
	ColumnNameCamel string
	ColumnNameSnake string //not used
	SqlType         string
	Seq             int
	Comment         string
	ColumnNameOut   string //dep: unclear what is the meaning
	GoTypeOut       string
	RoachTypeOut    string
	GoDefaultOut    string
	JavaTypeOut     string
	PBTypeOut       string
	StructTagOut    string
	IsNullAble      bool
	IsPrimary       bool
	IsUnique        bool
}

type ColumnType struct {
	SqlType    string
	GoType     string
	GoFlatType string
	JavaType   string
	PBType     string
}

type PrimaryKey struct {
	IsCompltive bool
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

// IndexColumn represents index column info.
type IndexColumn struct {
	SeqNo      int    // seq_no
	Cid        int    // cid
	ColumnName string // column_name
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
    s:= t.ColumnNameSnake + " " + t.RoachTypeOut + " "
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
