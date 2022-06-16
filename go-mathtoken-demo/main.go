package main

import (
	"fmt"
	"strings"

	"go-mathtoken-demo/mathtoken"
)

func main() {
	tokens, err := mathtoken.Parse("(10 + #device_login.#width.avg / 10")
	if err != nil {
		panic(err)
	}
	fmt.Println(tokens)
	values := make([]string, len(tokens))
	for i, t := range tokens {
		fmt.Println(t.Type, t.Value)
		values[i] = t.Value
	}

	fmt.Println("===> ", strings.Join(values, " "))
}
