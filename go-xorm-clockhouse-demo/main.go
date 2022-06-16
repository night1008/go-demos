package main

import (
	"fmt"

	"xorm.io/xorm"
)

func main() {
	engine, err := xorm.NewEngine("clickhouse", "tcp://localhost:9000")
	if err != nil {
		panic(err)
	}

	fmt.Println(engine)

	// results, err := engine.Query("select * from event_log")
	if err != nil {
		panic(err)
	}
}
