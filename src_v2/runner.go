package src_v2

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"regexp"
)

var IntRE = regexp.MustCompile(`^int(32|64)?$`)

func Run() {
	DB, err := sqlx.Connect("mysql", "flipper:12345678@tcp(192.168.162.115:3306)/flip_my?charset=utf8mb4")
	NoErr(err)
	DB.MapperFunc(func(s string) string { return s })
	DB = DB.Unsafe()

	//OutPutBuffer := &GenOut{}
	for _, db := range DATABASES {
		tables, err := mysql_loadTables(DB, db, "BASE TABLE")
		NoErr(err)
		_ = tables

		for _, table := range tables {
			rs := convNativeTableToOut(*table)
			PPJson(rs)
			fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
		}
		//OutPutBuffer.Tables = append(OutPutBuffer.Tables, tables...)
	}

	/*

		for _, table := range OutPutBuffer.Tables {
			if table.IsPG {
				continue
			}
			if table.NeedTrigger {
				//OutPutBuffer.TablesTriggers = append(OutPutBuffer.TablesTriggers, table)
			}
			if table.SinglePrimaryKey != nil {
				table.XPrimaryKeyGoType = table.SinglePrimaryKey.GoTypeOut
			}
		}

		// For Rust Filtering
		for _, table := range OutPutBuffer.Tables {
			if table.IsPG {
				continue
			}

			// Add Modifires
			for _, col := range table.Columns {
				col.WhereModifiersRust = col.GetModifiersRust()
				col.WhereInsModifiersRust = col.GetRustModifiersIns()
			}

			if table.HasPrimaryKey  {
				//OutPutBuffer.RustTables = append(OutPutBuffer.RustTables, table)
			}
		}

		//PertyPrint(OutPutBuffer.RustTables)
		PertyPrint(OutPutBuffer.Tables)
		rustBuild(OutPutBuffer)
		//goBuild(OutPutBuffer)
		//PertyPrint(OutPutBuffer.Tables)
	*/
}
