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
	DB, err := sqlx.Connect("mysql", "flipper:12345678@tcp(192.168.92.115:3306)/flip_my?charset=utf8mb4")
	NoErr(err)
	DB.MapperFunc(func(s string) string { return s })
	DB = DB.Unsafe()

	tables := []*OutTable{}

	for _, db := range DATABASES {
		nativeTables, err := mysql_loadTables(DB, db, "BASE TABLE")
		NoErr(err)
		_ = tables

		for _, table := range nativeTables {
			outTable := convNativeTableToOut(*table)
			//PPJson(outTable)
			fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
			tables = append(tables, outTable)
		}
	}
	PertyPrint(tables)

	// For Rust Filtering
	// Note: We add column modifiers in here in order to have a more shorter debug outputs in above of
	//	this line, originally we should add them in convNativeTableToOut().
	for _, outTable := range tables {
		// Add Modifiers
		for _, col := range outTable.Columns {
			col.WhereModifiersRust = col.GetModifiersRust()
			col.WhereInsModifiersRust = col.GetRustModifiersIns()
		}
	}

	//PPJson(OutPutBuffer)

	//PertyPrint(OutPutBuffer.Tables)

	OutPutBuffer := &GenOut{
		Tables:         nil,
		TablesFiltered: setFilteredTables(tables),
	}

	rustBuild(OutPutBuffer)
	//goBuild(OutPutBuffer)
	//PertyPrint(OutPutBuffer.Tables)

}

func setFilteredTables(tables []*OutTable) (res []*OutTable) {
	tablesFiltered := []*OutTable{}
	for _, t := range tables {
		// We can skip any tablesFiltered that we do not want in here. For now process all of them.
		// todo support multi primay keys
		if t.IsAutoIncr {
			//continue
		}
		tablesFiltered = append(tablesFiltered, t)
	}

	return tablesFiltered
}
