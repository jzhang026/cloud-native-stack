package main

import "fmt"

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	slice := arr[1:3]
	fmt.Printf("slice is %v", slice)

	fullSlice := arr[:]
	remove3rd := deleteItem(fullSlice, 3)
	fmt.Printf("remove3rd %v", remove3rd)

}

func deleteItem(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}
