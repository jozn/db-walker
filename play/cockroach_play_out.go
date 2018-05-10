package main

import (
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
    "fmt"
    "ms/db_walker/src"
    "ms/sun/shared/helper"
    "ms/db_walker/play/out"
)
var DB, err = sqlx.Connect("postgres", "postgresql://root@localhost:26257?sslmode=disable")

func main()  {
    src.RunCockRoach_Play_Dep()

    DB, err := sqlx.Connect("postgres", "postgresql://root@localhost:26257?sslmode=disable")
    fmt.Println(DB, err)
    s:= helper.SqlManyDollars(4,10,false)
    fmt.Println(s)
    _ = x.Account{}
}
