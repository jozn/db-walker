package src_v2

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// MyTables runs a custom query, returning results as Table.
func mysql_loadTables(db *sqlx.DB, schema string, relkind string) (res []*Table, err error) {
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

	fmt.Println("Mysql loader - load tables: ", res2)

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