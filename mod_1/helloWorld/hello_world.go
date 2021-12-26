package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	name := flag.String("name", "world", "input some name")
	flag.Parse()
	fmt.Println("OS is", os.Args)
	fmt.Println("input parameter is", name)
	fullString := fmt.Sprintf("hello %s from go", *name)
	fmt.Println(fullString)
}
