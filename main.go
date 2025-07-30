package main

import "fmt"

type Person struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age"`
}

func main() {
	fmt.Println("Hello World")
}
