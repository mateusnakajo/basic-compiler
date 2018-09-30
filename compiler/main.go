package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	switch {
	case len(args) > 1:
		fmt.Println("Usage: basic [script]")
	case len(args) == 1:
		fmt.Println("ok")
	case len(args) == 0:
		fmt.Println(">>")
	}
}
