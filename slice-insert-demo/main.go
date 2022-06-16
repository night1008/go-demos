package main

import "fmt"

func main() {
	names := []string{"foo", "bar", "x", "y"}
	fmt.Println(names[:len(names)])

	 

	names3 := names[:2]
	names3 = append(names3, "new")
	names3 = append(names3, names[2:]...)
	fmt.Println(names3)
}
