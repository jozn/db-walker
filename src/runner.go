package src

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var IntRE = regexp.MustCompile(`^int(32|64)?$`)

func Run() {
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
			OutPutBuffer.TablesTriggers = append(OutPutBuffer.TablesTriggers, table)
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
			OutPutBuffer.RustTables = append(OutPutBuffer.RustTables, table)
		}
	}

	PertyPrint(OutPutBuffer.RustTables)
	rustBuild(OutPutBuffer)
	//goBuild(OutPutBuffer)
	//PertyPrint(OutPutBuffer.Tables)

}

func addCockRoachTables(OutPutBuffer *GenOut) {
	DB, err := sqlx.Connect("postgres", "postgresql://root@localhost:26257?sslmode=disable")
	if err != nil {
		fmt.Println("cockroach connecting err: ", err)
		return
	}
	fmt.Println(DB, err)
	//on PG we must lowercase coulmns names unlike the Myql which is upper case
	DB.MapperFunc(func(s string) string { return strings.ToLower(s) })
	DB = DB.Unsafe()

	//OutPutBuffer := &GenOut{}
	for _, db := range DATABASES_COCKROACHE {
		tables, err := Cockroach_LoadTables(DB, db, "BASE TABLE")
		NoErr(err)
		OutPutBuffer.Tables = append(OutPutBuffer.Tables, tables...)
	}

	for _, table := range OutPutBuffer.Tables {
		if table.IsPG {
			table.Columns, _ = Cockroach_LoadTableColumns(DB, table.DataBase, table.TableName, table)
			table.Indexes, _ = Cockroach_TableIndexes(DB, table.DataBase, table.TableName, table)
		}
	}
}
