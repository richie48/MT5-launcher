package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const BASE_DIR string = "/mnt/c/Program Files/MetaTrader 5"
const CONFIG_FILE string = "config-example.json"
const PERMISSIONS fs.FileMode = 0755

type Account struct {
	Name     string
	Login    string
	Password string
	Server   string
	Path     string
}

func copyFile(source string, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	destinationFile, err := os.Create(destination)

	if err != nil {
		return err
	}

	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func recursiveCopy(sourceDirectory string, destinationDirectory string) {

	files, err := os.ReadDir(sourceDirectory)

	if err != nil {
		fmt.Println("error reading directory", sourceDirectory)
		return
	}

	for _, file := range files {
		sourcePath := filepath.Join(sourceDirectory, file.Name())
		destinationPath := filepath.Join(destinationDirectory, file.Name())

		if file.IsDir() {
			childSourceDirectory := sourceDirectory + "/" + file.Name()
			childDestinationDirectory := destinationDirectory + "/" + file.Name()

			err := os.Mkdir(childDestinationDirectory, PERMISSIONS)
			if err != nil {
				fmt.Println("error creating directory", childDestinationDirectory)
				return
			}

			recursiveCopy(childSourceDirectory, childDestinationDirectory)

		} else {

			if err := copyFile(sourcePath, destinationPath); err != nil {
				fmt.Println("Error copying files:", err)
			}
		}

	}
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

	recursiveCopy(BASE_DIR, instanceDirectory)

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
