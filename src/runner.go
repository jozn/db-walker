package src

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"ms/sun/helper"
	"regexp"
)

var IntRE = regexp.MustCompile(`^int(32|64)?$`)

var OutPutBuffer = &GenOut{
    PackageName: "x",
}

var EscapeColumnNames = false

func Run() {
	DB, err := sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/os?charset=utf8mb4")
	DB.MapperFunc(func(s string) string { return s })
	DB = DB.Unsafe()
	helper.NoErr(err)

	//OutPutBuffer := &GenOut{}
	OutPutBuffer.Tables, _ = My_LoadTables(DB, "os", "BASE TABLE")

	for _, table := range OutPutBuffer.Tables {
		table.Columns, _ = My_LoadTableColumns(DB, table.DataBase, table.TableName, table)
		table.Indexes, _ = MyTableIndexes(DB, table.DataBase, table.TableName, table)
	}

	build(OutPutBuffer)
	helper.PertyPrint(OutPutBuffer.Tables)

}
