package src_v2

import (
	"github.com/jmoiron/sqlx"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
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
		//OutPutBuffer.Tables = append(OutPutBuffer.Tables, tables...)
	}

	/*for _, table := range OutPutBuffer.Tables {
		table.Columns, _ = MySQL_LoadTableColumns(DB, table.DataBase, table.TableName, table)
		table.Indexes, _ = MySQL_TableIndexes(DB, table.DataBase, table.TableName, table)
	}

	// addCockRoachTables(OutPutBuffer)

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

func Run_old() {
	//DB, err := sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/sun?charset=utf8mb4")
	// DB, err := sqlx.Connect("mysql", "root:123456@tcp(37.152.187.1:3306)/twitter?charset=utf8mb4")
	// DB, err := sqlx.Connect("mysql", "root:12345678@tcp(130.185.120.132:3306)/twitter?charset=utf8mb4")
	DB, err := sqlx.Connect("mysql", "flipper:12345678@tcp(192.168.162.115:3306)/flip_my?charset=utf8mb4")
	NoErr(err)
	DB.MapperFunc(func(s string) string { return s })
	DB = DB.Unsafe()

	//OutPutBuffer := &GenOut{}
	for _, db := range DATABASES {
		tables, err := MySQL_LoadTables(DB, db, "BASE TABLE")
		NoErr(err)
		OutPutBuffer.Tables = append(OutPutBuffer.Tables, tables...)
	}

	for _, table := range OutPutBuffer.Tables {
		table.Columns, _ = MySQL_LoadTableColumns(DB, table.DataBase, table.TableName, table)
		table.Indexes, _ = MySQL_TableIndexes(DB, table.DataBase, table.TableName, table)
	}

	// addCockRoachTables(OutPutBuffer)

	for _, table := range OutPutBuffer.Tables {
		if table.IsPG {
			continue
		}
		if table.NeedTrigger {
			//OutPutBuffer.TablesTriggers = append(OutPutBuffer.TablesTriggers, table)
		}
		if table.PrimaryKey != nil {
			table.XPrimaryKeyGoType = table.PrimaryKey.GoTypeOut
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

		if table.HasPrimaryKey /*&& !table.IsCompositePrimaryKey*/ {
			//OutPutBuffer.RustTables = append(OutPutBuffer.RustTables, table)
		}
	}

	//PertyPrint(OutPutBuffer.RustTables)
	PertyPrint(OutPutBuffer.Tables)
	rustBuild(OutPutBuffer)
	//goBuild(OutPutBuffer)
	//PertyPrint(OutPutBuffer.Tables)

}

