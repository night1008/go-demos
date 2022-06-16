package main

import "fmt"

func main() {
	s1 := struct{}{}
	fmt.Println(fmt.Sprintf("%p", &s1))
	// 0x116ce80

	s2 := struct{}{}
	fmt.Println(fmt.Sprintf("%p", &s2))
	// 0x116ce80

	// different empty struct has same address

	// BTW
	b1 := true
	fmt.Println(fmt.Sprintf("%p", &b1))
	// 0xc00010c029

	b2 := true
	fmt.Println(fmt.Sprintf("%p", &b2))
	// 0xc00010c02a

	// same bool value has different address
}
