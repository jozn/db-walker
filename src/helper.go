package src

import (
	"fmt"
	"log"

	"github.com/kr/pretty"
)

func NoErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func PertyPrint(a interface{}) {
	fmt.Printf("%# v \n", pretty.Formatter(a))
}