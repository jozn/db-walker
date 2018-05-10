package main

import (
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
    "ms/db_walker/play/out"
    "github.com/jmoiron/sqlx"
    "fmt"
    "math/rand"
    "ms/db_walker/src"
)
var DB, err = sqlx.Connect("postgres", "postgresql://root@localhost:26257?sslmode=disable")

func main()  {
    src.RunCockRoach()

    DB, err := sqlx.Connect("postgres", "postgresql://root@localhost:26257?sslmode=disable")
    fmt.Println(DB, err)

    a := x.Account{
        Id: rand.Intn(100000),
        Balance:189.25,
    }

    err = a.Insert(DB)

    fmt.Println(err)

    f1()
    f2()
    f3()
    f4()


}

func f1()  {
    s,err := x.NewAccount_Selector().Id_GE(100).GetRows(DB)
    fmt.Println(err)
    fmt.Println(s)
}

func f2()  {
    s,err := x.NewAccount_Selector().Id_In([]int{100,101,92647}).Or().Id_Eq(25).GetRows(DB)
    fmt.Println(err)
    fmt.Println(s)
}

func f3()  {
    s,err := x.NewAccount_Deleter().Id_In([]int{100,101,92647}).Or().Id_Eq(25).Delete(DB)
    fmt.Println(err)
    fmt.Println(err)
    fmt.Println(s)
}

func f4()  {
    s,err := x.NewAccount_Updater().Id_Increment(1).Id_In([]int{100,101,92647}).Or().Id_Eq(25).Update(DB)
    fmt.Println(err)
    fmt.Println(s)
}
