package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
)

const BASE_DIR string = "/mnt/c/Program Files/MetaTrader 5/"
const CONFIG_FILE string = "config-example.json"
const PERMISSIONS fs.FileMode = 0755

type Account struct {
	Name     string
	Login    string
	Password string
	Server   string
	Path     string
}

func createInstance(account Account) {
	var instanceDirectory string = "MT5_" + account.Name

	if _, err := os.Stat(instanceDirectory); !os.IsNotExist(err) {
		fmt.Println(instanceDirectory, "already exist, skipping!")
		return
	}

	fmt.Println("creating instance ", instanceDirectory)

	if err := os.MkdirAll(instanceDirectory, PERMISSIONS); err != nil {
		fmt.Println("error creating directory", instanceDirectory)
	}

	files, err := os.ReadDir(BASE_DIR)

	if err != nil {
		fmt.Println("error reading directory", BASE_DIR)
		return
	}

	for _, file := range files {
		// TODO: copy files to destination

	}

}

func main() {
	jsonFile, err := os.Open(CONFIG_FILE)

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
		createInstance(account)
	}
}
