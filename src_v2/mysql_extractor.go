package src_v2

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"strings"
)

// MyTables runs a custom query, returning results as Table.
func mysql_loadTables(db *sqlx.DB, schema string, relkind string) (res []*NativeTable, err error) {
	// sql query
	const sqlstr = `SELECT * ` +
		`FROM information_schema.tables ` +
		`WHERE table_schema = ? AND table_type = ?`

	// run query
	XOLogDebug(sqlstr, schema, relkind)

	var tabels = []struct {
		TABLE_NAME     string
		TABLE_TYPE string
		ENGINE string
		TABLE_COMMENT string
		// Note: This filed just is the counter for auto_increment, in newly created tables this is null even if
		//	the table has auto_increment column, use EXTRA column in inforamation_schema
		AUTO_INCREMENT sql.NullInt64
	}{}
	err = db.Unsafe().Select(&tabels, sqlstr, schema, relkind)
	NoErr(err)

	//fmt.Println("Mysql loader - load tables: ", tabels)
	//PPJson(tabels)

	for _, table := range tabels {
		// Load Columns
		cols, err := mysql_loadTableColumns(db,schema,table.TABLE_NAME)
		NoErr(err)

		nt := &NativeTable{
			TableName:             table.TABLE_NAME,
			Columns:               cols,
			HasPrimaryKey:         false,
			IsCompositePrimaryKey: false,
			SinglePrimaryKey:      nil,
			PrimaryKeys:           nil,
			DataBase:              schema,
			Comment:               table.TABLE_COMMENT,
			IsAutoIncrement:       false,
			Indexes:               nil,
		}

		// Set table data info about columns info
		primaryCols :=  []*NativeColumn{}
		for _, col := range cols {
			if col.IsPrimary {
				nt.HasPrimaryKey = true

				primaryCols =  append(primaryCols, col)
			}

			if col.IsAutoIncrement {
				nt.IsAutoIncrement = true
			}
		}

		// Primary columns set
		if len(primaryCols) == 1 {
			nt.SinglePrimaryKey = primaryCols[0]
		} else if len(primaryCols) >= 2 {
			nt.PrimaryKeys = primaryCols
		}

		// Load Indexes

		indxs, err := mysql_TableIndexes(db,schema,table.TABLE_NAME, &Table{})
		NoErr(err)
		PPJson(indxs)

		res = append(res, nt)
	}

	//PertyPrint(res[0])
	//PPJson(res)

	return res, nil
}

func mysql_loadTableColumns(db *sqlx.DB, schema string, tableName string) (res []*NativeColumn, err error) {
	// Notes:
	//	+ "MUL" in COLUMN_KEY: means it has an not unique index
	//	+ "UNI" in COLUMN_KEY: means it has an unique index
	//	+ "PRI" in COLUMN_KEY: multi columns could have this in case of compund index

	var colRows = []struct {
		ORDINAL_POSITION int // Starts form 0
		COLUMN_NAME      string	// ex: "channel_msg"
		DATA_TYPE        string // simple: "varchar" "bigint" "text" "intt" >> no size limit and unsigned description in here
		IS_NULLABLE      string // 'YES' 'NO'
		COLUMN_DEFAULT   sql.NullString // null or string like '0'
		COLUMN_TYPE      string // "varchar(50)" "bigint unsigned" " text" "blob" >> with size limit and unsigned description
		COLUMN_KEY       string // if == 'PRI' then is the primiry key -- not neccoery auto_incer -- "PRI", "UNI", "MUL", ""
		EXTRA            string // if == 'auto_increment' then this is the auto incerment -- not neccoery primiry key
		COLUMN_COMMENT   string
	}{}
	// sql query
	const sqlstr = `SELECT * ` +
		`FROM information_schema.columns ` +
		`WHERE table_schema = ? AND table_name = ? ` +
		`ORDER BY ordinal_position ASC`

	// run query
	XOLogDebug(sqlstr, schema, tableName)

	err = db.Unsafe().Select(&colRows, sqlstr, schema, tableName)
	NoErr(err)
	//fmt.Println("+++++ Mysql loader - load tables: ", colRows)
	for _, colRow := range colRows {
		//if this coulmn is auto_incermnt but not primiry this means: this table has one auto Seq columns
		//so skip it from our entire genrated paradigram and make the table
		// Updated in Rust version: we do not support this functionality anymore > commented
		if colRow.EXTRA == "auto_increment" && colRow.COLUMN_KEY != "PRI" {
			//table.IsAutoIncrement = false
			//continue // Skip this table in generated code
		}

		nullable := false
		switch colRow.IS_NULLABLE {
		case "YES":
			nullable = true
		case "NO":
			nullable = false
		}

		col := &NativeColumn{
			ColumnName:      colRow.COLUMN_NAME,
			SqlType:         colRow.DATA_TYPE,
			SqlTypeFull:     colRow.COLUMN_TYPE,
			Ordinal:         colRow.ORDINAL_POSITION,
			Comment:         colRow.COLUMN_COMMENT,
			IsNullAble:      nullable,
			// Set below
			IsPrimary:       false,
			IsUnique:        false,
			IsAutoIncrement: false,
		}

		if colRow.COLUMN_KEY == "UNI" {
			col.IsUnique = true
		}

		if colRow.COLUMN_KEY == "PRI" {
			col.IsPrimary = true
		}

		if colRow.EXTRA == "auto_increment" {
			col.IsAutoIncrement = true
		}

		res = append(res, col)
	}

	return res, nil
}

func mysql_TableIndexes(db *sqlx.DB, schema string, tableName string, table *Table) (res []*Index, err error) {
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
		//i.FuncNameOut = GoIndexName(i, table)
		//i.RustFuncName = RustIndexName(i, table)
		res = append(res, i)
	}

	return res, nil
}
