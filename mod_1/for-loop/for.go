package main

import "fmt"

func main() {
	for i := 0; i < 9; i++ {
		fmt.Println(i)
	}

	fullString := "hello world!"

	for i, char := range fullString {
		fmt.Println(i, string(char))
	}
}
