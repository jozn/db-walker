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

	OutPutBuffer := &GenOut{}
	for _, db := range DATABASES {
		tables, err := mysql_loadTables(DB, db, "BASE TABLE")
		NoErr(err)
		_ = tables

		for _, table := range tables {
			outTable := convNativeTableToOut(*table)
			//PPJson(outTable)
			fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			OutPutBuffer.Tables = append(OutPutBuffer.Tables, outTable)
		}
	}
	PertyPrint(OutPutBuffer.Tables)

	// For Rust Filtering
	// Note: We add column modifiers in here in order to have a more shorter debug outputs in above of
	//	this line, originally we should add them in convNativeTableToOut().
	for _, outTable := range OutPutBuffer.Tables {
		// Add Modifiers
		for _, col := range outTable.Columns {
			col.WhereModifiersRust = col.GetModifiersRust()
			col.WhereInsModifiersRust = col.GetRustModifiersIns()
		}
	}

	//PPJson(OutPutBuffer)

	//PertyPrint(OutPutBuffer.Tables)

	setFilteredTables(OutPutBuffer)

	rustBuild(OutPutBuffer)
	//goBuild(OutPutBuffer)
	//PertyPrint(OutPutBuffer.Tables)

}

func setFilteredTables(gen *GenOut) {
	tables := []*OutTable{}
	for _, t := range gen.Tables {
		// We can skip any tables that we do not want in here. For now process all of them.
		// todo support multi primay keys
		if t.SinglePrimaryKey == nil {
			//continue
		}
		tables = append(tables, t)
	}

	gen.TablesFiltered = tables
}
