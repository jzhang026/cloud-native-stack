package main

import "fmt"

func main() {
	arr := [5]string{"I", "am", "stupid", "and", "weak"}
	fmt.Printf("before %v\n", arr)

	for idx, _ := range arr {
		switch idx {
		case 2:
			arr[idx] = "smart"
		case 4:
			arr[idx] = "strong"
		}
	}

	fmt.Printf("after %v\n", arr)
}
