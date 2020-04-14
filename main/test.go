package main

import (
	"fmt"
	"strconv"
)

func main() {
	testStr := make([]string, 20)
	fmt.Println(testStr[11])
	for i:= range testStr {
		testStr[i] = strconv.Itoa(i)
	}
	fmt.Println(testStr)
	testStr=nil
	fmt.Println(testStr)
}
