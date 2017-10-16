package src

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"ms/sun/helper"
)

func Run() {
	DB, err := sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/os?charset=utf8mb4")
	DB.MapperFunc(func(s string) string { return s })
	DB = DB.Unsafe()
	helper.NoErr(err)

	tables, _ := My_LoadTables(DB, "os", "BASE TABLE")

	for _, t := range tables {
		t.Columns, _ = My_LoadTableColumns(DB, t.DataBase, t.TableName)
		t.Indexes, _ = MyTableIndexes(DB, t.DataBase, t.TableName, t)
	}

	helper.PertyPrint(tables)

}
