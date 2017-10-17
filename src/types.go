package src

type GenOut struct {
    PackageName          string
	Tables               []*Table
	GeneratedPb          string
	GeneratedPbConverter string
}

//dep
type DataBase struct {
	Tables []Table
}

type Table struct {
	TableName       string
	Columns         []*Column
	HasPrimaryKey   bool
	PrimaryKey      *Column
	DataBase        string //or schema in PG or tablesapce in cassandra
	Seq             int
	Comment         string
	IsAutoIncrement bool
	Indexes         []*Index
	TableNameOut    string //with table "`ms`.`post`"
	TableNameGo     string
	TableNameJava   string
	TableNamePB     string
	ShortName       string
}

type Column struct {
	ColumnName    string
	SqlType       string
	Seq           int
	Comment       string
	ColumnNameOut string
	GoTypeOut     string
	GoDefaultOut  string
	JavaTypeOut   string
	PBTypeOut     string
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
	IndexName string // index_name
	IsUnique  bool   // is_unique
	IsPrimary bool   // is_primary
	SeqNo     int    // seq_no
	Columns   []*Column
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
