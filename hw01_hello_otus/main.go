package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	sample := "Hello, OTUS!"
	fmt.Print(stringutil.Reverse(sample))
}
