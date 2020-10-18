package main

import (
	"fmt"
	"strings"
)

func main() {
	sampleString := "https://me.titan.in/"
	sampleString2 := "/this-is-a-path/to"
	fmt.Println(sampleString, " ", strings.Trim(sampleString, "/"))
	fmt.Println(sampleString2, " ", strings.Trim(sampleString2, "/"))
}
