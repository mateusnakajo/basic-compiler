package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	args := os.Args[1:]
	switch {
	case len(args) > 1:
		fmt.Println("Usage: basic [script]")
	case len(args) == 1:
		dat, err := ioutil.ReadFile(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Print(string(dat))
	case len(args) == 0:
		fmt.Println(">>")
	}
}
