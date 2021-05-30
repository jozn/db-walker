package src_v2

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Notes:
//	+ MySQL does not set unique to true for primary keys even though they are.

func mysql_loadTables(db *sqlx.DB, schema string, relkind string) (res []*NativeTable, err error) {
	// sql query
	const sqlstr = `SELECT * ` +
		`FROM information_schema.tables ` +
		`WHERE table_schema = ? AND table_type = ?`

	// run query
	XOLogDebug(sqlstr, schema, relkind)

	var tabels = []struct {
		TABLE_NAME    string
		TABLE_TYPE    string
		ENGINE        string
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
		cols, err := mysql_loadTableColumns(db, schema, table.TABLE_NAME)
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
		primaryCols := []*NativeColumn{}
		for _, col := range cols {
			if col.IsPrimary {
				nt.HasPrimaryKey = true

				primaryCols = append(primaryCols, col)
			}

			if col.IsAutoIncrement {
				nt.IsAutoIncrement = true
			}
		}

		// IsPrimary columns set
		if len(primaryCols) == 1 {
			nt.SinglePrimaryKey = primaryCols[0]
		} else if len(primaryCols) >= 2 {
			nt.PrimaryKeys = primaryCols
		}

		// Load Indexes

		//indxs, err := mysql_TableIndexes_old(db,schema,table.TABLE_NAME, &Table{})
		indxs, err := mysql_loadIndexs(db, schema, table.TABLE_NAME, nt)
		nt.Indexes = indxs
		NoErr(err)
		//fmt.Println("&&&&&&&&&&&&&&&&&&&&& indes")
		//PPJson(indxs)

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
		ORDINAL_POSITION int            // Starts form 0
		COLUMN_NAME      string         // ex: "channel_msg"
		DATA_TYPE        string         // simple: "varchar" "bigint" "text" "intt" >> no size limit and unsigned description in here
		IS_NULLABLE      string         // 'YES' 'NO'
		COLUMN_DEFAULT   sql.NullString // null or string like '0'
		COLUMN_TYPE      string         // "varchar(50)" "bigint unsigned" " text" "blob" >> with size limit and unsigned description
		COLUMN_KEY       string         // if == 'PRI' then is the primiry key -- not neccoery auto_incer -- "PRI", "UNI", "MUL", ""
		EXTRA            string         // if == 'auto_increment' then this is the auto incerment -- not neccoery primiry key
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
			ColumnName:  colRow.COLUMN_NAME,
			SqlType:     colRow.DATA_TYPE,
			SqlTypeFull: colRow.COLUMN_TYPE,
			Ordinal:     colRow.ORDINAL_POSITION,
			Comment:     colRow.COLUMN_COMMENT,
			IsNullAble:  nullable,
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

// This func does not change the (table *NativeTable) fields.
func mysql_loadIndexs(db *sqlx.DB, schema string, tableName string, table *NativeTable) (res []*NativeIndex, err error) {
	// Notes:
	// + PRIMARY key is just a UNIQUE NOT NULL constraint (https://dev.mysql.com/doc/refman/8.0/en/innodb-index-types.html)
	// + MySQL seems to not allow change the name of PRIMARY Index

	type tableRow struct {
		NON_UNIQUE int // 0 or 1 -- 1 is being set just for none-primary NOT unique types (multi rows)
		//INDEX_SCHEMA string // Name of table
		INDEX_NAME    string // 'PRIMARY' or other name of index
		SEQ_IN_INDEX  int    // Strarts from 1
		COLUMN_NAME   string
		NULLABLE      string // "YES" or ""
		INDEX_TYPE    string // "BTREE" "HASH"
		INDEX_COMMENT string
	}

	var colRows = []*tableRow{}
	// sql query
	const sqlstr = `SELECT * ` +
		`FROM information_schema.statistics ` +
		`WHERE table_schema = ? AND table_name = ? ` +
		`ORDER BY seq_in_index ASC `
	// run query
	XOLogDebug(sqlstr, schema, tableName)

	err = db.Unsafe().Select(&colRows, sqlstr, schema, tableName)
	NoErr(err)
	//PPJson(colRows)
	mp := make(map[string][]*tableRow)
	for _, colRow := range colRows {
		mp[colRow.INDEX_NAME] = append(mp[colRow.INDEX_NAME], colRow)
	}

	for idxName, idxCols := range mp {
		idx1 := idxCols[0]

		isUnique := true
		if idx1.NON_UNIQUE == 1 {
			isUnique = false
		}

		ic := &NativeIndex{
			IndexName: idxName,
			IsUnique:  isUnique,
			IsPrimary: idxName == "PRIMARY",
			ColNum:    len(idxCols),
			Comment:   idx1.INDEX_COMMENT,
			Columns:   nil,
			Table:     nil,
		}

		for _, col := range idxCols {
			for _, nt := range table.Columns {
				if nt.ColumnName == col.COLUMN_NAME {
					ic.Columns = append(ic.Columns, nt)
				}
			}
		}

		res = append(res, ic)
	}
	//PPJson(mp)

	return res, nil
}

var XOLogDebug = func(s string, o ...interface{}) {
	if true {
		fmt.Println(s, o)
	}
}
