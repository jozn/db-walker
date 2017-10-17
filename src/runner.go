package src

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"ms/sun/helper"
	"regexp"
)

var IntRE = regexp.MustCompile(`^int(32|64)?$`)

var OutPutBuffer = OutPut{}

func Run() {
	DB, err := sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/os?charset=utf8mb4")
	DB.MapperFunc(func(s string) string { return s })
	DB = DB.Unsafe()
	helper.NoErr(err)

	tables, _ := My_LoadTables(DB, "os", "BASE TABLE")

	for _, table := range tables {
		table.Columns, _ = My_LoadTableColumns(DB, table.DataBase, table.TableName, table)
		table.Indexes, _ = MyTableIndexes(DB, table.DataBase, table.TableName, table)
	}

	helper.PertyPrint(tables)

}
