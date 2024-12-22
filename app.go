package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	// ask for user input
	var name string
	fmt.Print("What is your name? ")
	fmt.Scanln(&name)
	fmt.Println("Hello", name)

}
