package src_v2

import (
	"fmt"
	"log"

	"github.com/hokaccha/go-prettyjson"
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

func PertyPrint2(a interface{}) {
	fmt.Printf("%# v \n", ToJsonPerety(a))
}
func ToJsonPerety(structoo interface{}) string {
	bts, _ := prettyjson.Marshal(structoo)
	return string(bts)
}

func PPJson(structObj interface{})  {
	fmt.Println(ToJsonPerety(structObj))
	fmt.Println("================================")
}