package src

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"fmt"
	"github.com/jmoiron/sqlx"
	"ms/sun/shared/helper"
	"regexp"
	"strings"
)

var IntRE = regexp.MustCompile(`^int(32|64)?$`)

func Run() {
	DB, err := sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/sun?charset=utf8mb4")
	DB.MapperFunc(func(s string) string { return s })
	DB = DB.Unsafe()
	helper.NoErr(err)

	//OutPutBuffer := &GenOut{}
	for _, db := range DATABASES {
		tables, err := My_LoadTables(DB, db, "BASE TABLE")
		helper.NoErr(err)
		OutPutBuffer.Tables = append(OutPutBuffer.Tables, tables...)
	}

	for _, table := range OutPutBuffer.Tables {
		table.Columns, _ = My_LoadTableColumns(DB, table.DataBase, table.TableName, table)
		table.Indexes, _ = MyTableIndexes(DB, table.DataBase, table.TableName, table)
	}

    addCockRoachTables(OutPutBuffer)

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

    helper.PertyPrint(OutPutBuffer.Tables)
	build(OutPutBuffer)
	//helper.PertyPrint(OutPutBuffer.Tables)

}

func addCockRoachTables(OutPutBuffer *GenOut) {
	DB, err := sqlx.Connect("postgres", "postgresql://root@localhost:26257?sslmode=disable")
	fmt.Println(DB, err)
	//on PG we must lowercase coulmns names unlike the Myql which is upper case
	DB.MapperFunc(func(s string) string { return strings.ToLower(s) })
	DB = DB.Unsafe()
	if err != nil {
		fmt.Println("cockroach connecting err: ", err)
	}
	//OutPutBuffer := &GenOut{}
	for _, db := range DATABASES_COCKROACHE {
		tables, err := Roach_LoadTables(DB, db, "BASE TABLE")
		helper.NoErr(err)
		OutPutBuffer.Tables = append(OutPutBuffer.Tables, tables...)
	}

	for _, table := range OutPutBuffer.Tables {
	    if table.IsPG {
            table.Columns, _ = Roach_LoadTableColumns(DB, table.DataBase, table.TableName, table)
            table.Indexes, _ = RoachTableIndexes(DB, table.DataBase, table.TableName, table)
        }
	}
}

func RunCockRoach_Play_Dep() {
	//DB, err := sqlx.Connect("mysql", "root:123456@tcp(localhost:3306)/sun?charset=utf8mb4")
	//DB, err := sqlx.Connect("postgres", "user=root dbname=bank sslmode=disable")
	DB, err := sqlx.Connect("postgres", "postgresql://root@localhost:26257?sslmode=disable")
	fmt.Println(DB, err)
	//on PG we must lowercase coulmns names unlike the Myql which is upper case
	DB.MapperFunc(func(s string) string { return strings.ToLower(s) })
	DB = DB.Unsafe()
	helper.NoErr(err)

	//OutPutBuffer := &GenOut{}
	for _, db := range DATABASES_COCKROACHE {
		tables, err := Roach_LoadTables(DB, db, "BASE TABLE")
		helper.NoErr(err)
		OutPutBuffer.Tables = append(OutPutBuffer.Tables, tables...)
	}

	for _, table := range OutPutBuffer.Tables {
		table.Columns, _ = Roach_LoadTableColumns(DB, table.DataBase, table.TableName, table)
		table.Indexes, _ = RoachTableIndexes(DB, table.DataBase, table.TableName, table)
	}

	for _, table := range OutPutBuffer.Tables {
		if table.NeedTrigger {
			OutPutBuffer.TablesTriggers = append(OutPutBuffer.TablesTriggers, table)
		}
		if table.PrimaryKey != nil {
			table.XPrimaryKeyGoType = table.PrimaryKey.GoTypeOut
		}
	}

	build(OutPutBuffer)
	helper.PertyPrint(OutPutBuffer.Tables)

}
