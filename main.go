package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Account struct {
	Name     string
	Login    string
	Password string
	Server   string
	Path     string
}

func main() {
	jsonFile, err := os.Open("config-example.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("opened json file", jsonFile)

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		fmt.Println(err)
	}

	var accounts []Account
	if err := json.Unmarshal(byteValue, &accounts); err != nil {
		fmt.Println(err)
	}

	for _, account := range accounts {
		fmt.Printf("account name: " + account.Name)
	}
}
