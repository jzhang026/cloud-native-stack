package main

import "reflect"

type MyStruct struct {
	Name string `json:"name"`
}

func main() {
	mt := MyStruct{Name: "test"}
	myType := reflect.TypeOf(mt)
	name := myType.Field(0)
	tag := name.Tag.Get("json")
	print(tag)
}
