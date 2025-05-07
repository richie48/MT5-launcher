package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	jsonFile, err := os.Open("config.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	fmt.Println("opened json file", jsonFile)
}
