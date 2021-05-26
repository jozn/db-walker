package src_v2

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

var XOLogDebug = func(s string, o ...interface{}) {
	if false {
		fmt.Println("=SQL=> ", s, o)
	}
}

// MyTables runs a custom query, returning results as Table.
func MySQL_LoadTables(db *sqlx.DB, schema string, relkind string) (res []*Table, err error) {
	// sql query
	const sqlstr = `SELECT * ` +
		`FROM information_schema.tables ` +
		`WHERE table_schema = ? AND table_type = ?`

	// run query
	XOLogDebug(sqlstr, schema, relkind)

	var res2 = []struct {
		TABLE_NAME     string //`json:"rec_created_by"  db:"TABLE_NAME"`
		// Note: This filed just is the counter for auto_increment, in newly created tables this is null even if
		//	the table has auto_increment column, use EXTRA column in inforamation_schema
		AUTO_INCREMENT sql.NullInt64
	}{}
	err = db.Unsafe().Select(&res2, sqlstr, schema, relkind)
	NoErr(err)

	//fmt.Println("Mysql loader - load tables: ", res2)

	for i, r := range res2 {
		t := &Table{
			TableName:      r.TABLE_NAME,
			TableSchemeOut: fmt.Sprintf("%s.%s", schema, r.TABLE_NAME), //fmt.Sprintf("`%s`.`%s`", schema, r.TABLE_NAME),
			TableNameSql:   r.TABLE_NAME,                               //fmt.Sprintf("`%s`.`%s`", schema, r.TABLE_NAME),
			DataBase:       schema,
			Seq:            i,
			TableNameGo:    SingularizeIdentifier(r.TABLE_NAME), //,SnakeToCamel(r.TABLE_NAME),
			TableNameJava:  SingularizeIdentifier(r.TABLE_NAME), //SnakeToCamel(r.TABLE_NAME),
			TableNameRust:  SingularizeIdentifier(r.TABLE_NAME), //SnakeToCamel(r.TABLE_NAME),
			//TableNamePB:    "PB_" + SingularizeIdentifier(r.TABLE_NAME), //SnakeToCamel(r.TABLE_NAME),
			TableNamePB: "" + SingularizeIdentifier(r.TABLE_NAME), //SnakeToCamel(r.TABLE_NAME),
			ShortName:   shortname(r.TABLE_NAME, "err", "res", "sqlstr", "db", "XOLog"),
			NeedTrigger: needTriggerTable(r.TABLE_NAME),
			IsMysql:     true,
			IsPG:        false,
			Dollar:      "?",
		}
		if r.AUTO_INCREMENT.Valid {
			t.IsAutoIncrement = true
		}
		if t.NeedTrigger {

		}
		res = append(res, t)
	}
	//PertyPrint(res)

	return res, nil
}

// MySQL_LoadTableColumns runs a custom query, returning results as Column.
func MySQL_LoadTableColumns(db *sqlx.DB, schema string, tableName string, table *Table) (res []*Column, err error) {
	var rows = []struct {
		ORDINAL_POSITION int
		COLUMN_NAME      string
		DATA_TYPE        string
		IS_NULLABLE      string //'YES'
		COLUMN_DEFAULT   sql.NullString
		COLUMN_TYPE      string
		COLUMN_KEY       string //if == 'PRI' then is the primiry key -- not neccoery auto_incer
		EXTRA            string //if == 'auto_increment' then this is the auto incerment -- not neccoery primiry key
		COLUMN_COMMENT   string
	}{}
	// sql query
	const sqlstr = `SELECT * ` +
		`FROM information_schema.columns ` +
		`WHERE table_schema = ? AND table_name = ? ` +
		`ORDER BY ordinal_position ASC`

	// run query
	XOLogDebug(sqlstr, schema, tableName)

	err = db.Unsafe().Select(&rows, sqlstr, schema, tableName)
	NoErr(err)
	//fmt.Println("Mysql loader - load tables: ", rows)
	for _, r := range rows {
		//if this coulmn is auto_incermnt but not primiry this means: this table has one auto Seq columns
		//so skip it from our entire genrated paradigram and make the table
		// Updated in Rust version: we do not support this functionality anymore > commented
		if strings.ToLower(r.EXTRA) == "auto_increment" && strings.ToUpper(r.COLUMN_KEY) != "PRI" {
			//table.IsAutoIncrement = false
			//continue // Skip this table in generated code
		}

		nullable := false
		switch strings.ToUpper(r.IS_NULLABLE) {
		case "YES":
			nullable = true
		case "NO":
			nullable = false
		}

		//fmt.Println(r)
		typRs, typOrgRs, _ := sqlTypesToRustType(r.DATA_TYPE)
		col := &Column{
			ColumnName:      r.COLUMN_NAME,
			ColumnNameCamel: SnakeToCamel(r.COLUMN_NAME),
			ColumnNameSnake: ToSnake(r.COLUMN_NAME),
			Seq:             r.ORDINAL_POSITION,
			Comment:         r.COLUMN_COMMENT,
			ColumnNameOut:   r.COLUMN_NAME,
			SqlType:         r.COLUMN_TYPE,
			SqlTypeStrip:    r.DATA_TYPE,
			StructTagOut:    fmt.Sprintf("`db:\"%s\"`", r.COLUMN_NAME),
			IsNullAble:      nullable,
			// Rust
			ColumnNameRust: r.COLUMN_NAME,
			RustTypeOut:    typRs,
			TypeRustBorrow: typOrgRs,
		}

		if strings.ToUpper(r.COLUMN_KEY) == "PRI" {
			col.IsPrimary = true
			if table.HasPrimaryKey {
				table.IsCompositePrimaryKey = true
			}
			table.HasPrimaryKey = true
			table.PrimaryKey = col
			table.PrimaryKeys = append(table.PrimaryKeys, col)
		}
		if strings.ToUpper(r.COLUMN_KEY) == "UNI" {
			col.IsUnique = true
		}
		if strings.ToLower(r.EXTRA) == "auto_increment" {
			col.IsAutoIncrement = true
			table.IsAutoIncrement = true
		}
		//fmt.Println("Mysql loader - load tables: ))))))) ", col)
		res = append(res, col)
	}

	return res, nil
}

// MySQL_TableIndexes runs a custom query, returning results as Index.
func MySQL_TableIndexes(db *sqlx.DB, schema string, tableName string, table *Table) (res []*Index, err error) {
	// sql query
	var rows = []struct {
		INDEX_NAME string
		IS_UNIQUE  bool
	}{}

	const sqlstr = `SELECT ` +
		`DISTINCT INDEX_NAME, ` +
		`NOT non_unique AS IS_UNIQUE ` +
		`FROM information_schema.statistics ` +
		//`WHERE index_name <> 'PRIMARY' AND index_schema = ? AND table_name = ?`
		`WHERE index_schema = ? AND table_name = ? AND INDEX_NAME not like '%skip%' `

	XOLogDebug(sqlstr, schema, tableName)
	err = db.Select(&rows, sqlstr, schema, tableName)
	if err != nil {
		NoErr(err)
		return
	}

	for _, r := range rows {
		i := &Index{
			IndexName: r.INDEX_NAME,
			IsUnique:  r.IS_UNIQUE,
			//FuncNameOut: "Get" + table.TableNameGo + "By" + r.INDEX_NAME,
		}
		if strings.ToUpper(r.INDEX_NAME) == "PRIMARY" {
			i.IsPrimary = true
		}

		rs := []struct {
			SEQ_IN_INDEX int
			COLUMN_NAME  string
		}{}
		// sql query
		const sqlstr = `SELECT * ` +
			//`seq_in_index, ` + //starts from 1
			//`column_name ` +
			`FROM information_schema.statistics ` +
			`WHERE index_schema = ? AND table_name = ? AND index_name = ? ` +
			`ORDER BY seq_in_index`

		XOLogDebug(sqlstr, schema, tableName, i.IndexName)
		err = db.Unsafe().Select(&rs, sqlstr, schema, tableName, i.IndexName)
		if err != nil {
			NoErr(err)
			return
		}

		for _, c := range rs {
			i.Columns = append(i.Columns, table.GetColumnByName(c.COLUMN_NAME))
		}
		i.FuncNameOut = GoIndexName(i, table)
		i.RustFuncName = RustIndexName(i, table)
		res = append(res, i)
	}

	return res, nil
}

func GoIndexName(index *Index, table *Table) string {
	name := ""
	//PertyPrint(table)
	//PertyPrint(index)
	if len(index.Columns) == 1 {
		//name = "Get" + table.TableNameGo + "By" + index.Columns[0].ColumnName
		name = "" + table.TableNameGo + "By" + index.Columns[0].ColumnNameCamel
	} else {
		arr := []string{}
		for _, col := range table.Columns {
			arr = append(arr, col.ColumnNameCamel)
		}
		//name = "Get" + table.TableNameGo + "By" + strings.Join(arr, "And")
		name = "" + table.TableNameGo + "By" + strings.Join(arr, "And")
	}

	return name
}

// In rust version we use only index names for building fun names
// format: tweet_media_[[media_type]]_IDX > only [[ ]] is used for nameing
func RustIndexName(index *Index, table *Table) string {
	name := ""
	if index.IsPrimary {
		name = "get_" + table.TableName + ""
		return name
	}
	stripName := strings.Replace(index.IndexName, table.TableName+"_", "", -1)
	stripName = strings.Replace(stripName, "_IDX", "", -1)

	name = "get_" + table.TableName + "_" + stripName

	return name
}
