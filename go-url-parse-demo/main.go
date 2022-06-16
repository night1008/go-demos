package main

import (
	"fmt"
	"net/url"
	"strconv"
)

func main() {
	url, err := url.Parse("tcp://localhost:9000?compress=0")
	if err != nil {
		panic(err)
	}
	query := url.Query()
	fmt.Println("===> query.Get(\"compress\") ", query.Get("compress"))
	if v, err := strconv.ParseBool(query.Get("compress")); err != nil {
		panic(err)
	} else {
		fmt.Println("==> query compress", v)
	}
}
