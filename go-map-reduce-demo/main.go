package main

import (
	"fmt"
	"strings"
)

func mapStrToStr(strs []string, fn func(string) string) []string {
	var newStrs []string
	for _, str := range strs {
		newStrs = append(newStrs, fn(str))
	}
	return newStrs
}

func main() {
	strs := []string{"hello", "world", "foo"}
	fmt.Println(mapStrToStr(strs, func(s string) string {
		return strings.ToUpper(s)
	}))
}

// 控制逻辑 与 业务逻辑 分离
