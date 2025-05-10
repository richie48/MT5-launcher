package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const ConfigFile string = "config.json"
const FolderPrefix string = "MT5_"
const Permissions fs.FileMode = 0755

// Account is used to represent individual config for the array in config.json
type Account struct {
	Name     string
	Login    string
	Password string
	Server   string
	Path     string
}

// copyFile copies a file from 'source' path to the 'destination' path provided
func copyFile(source string, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		fmt.Errorf("error=", err)
		return err
	}

	defer sourceFile.Close()

	destinationFile, err := os.Create(destination)
	if err != nil {
		fmt.Errorf("error=", err)
		return err
	}

	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// recursiveCopy recursively copies all the files in the 'sourceDirectory' to
// 'destinationDirectory'. Return an error if anything goes wrong'
func recursiveCopy(sourceDirectory string, destinationDirectory string) error {

	files, err := os.ReadDir(sourceDirectory)
	if err != nil {
		fmt.Errorf("error reading directory", sourceDirectory)
		return err
	}

	var fileCopyError error

	for _, file := range files {
		sourcePath := filepath.Join(sourceDirectory, file.Name())
		destinationPath := filepath.Join(destinationDirectory, file.Name())

		if file.IsDir() {
			childSourceDirectory := sourceDirectory + "/" + file.Name()
			childDestinationDirectory := destinationDirectory + "/" + file.Name()

			err := os.Mkdir(childDestinationDirectory, Permissions)
			if err != nil {
				fmt.Errorf("error creating directory", childDestinationDirectory)
				return err
			}

			fileCopyError = recursiveCopy(childSourceDirectory, childDestinationDirectory)

		} else {

			err := copyFile(sourcePath, destinationPath)
			if err != nil {
				fmt.Errorf("Error copying files", err)
				return err
			}
		}

	}

	return fileCopyError
}

// CreateInstance attempts to create a directory. Its name prefixed with 'FolderPrefix'
// and some details in 'account'. The aim of this directory is to replicate everything
// fould in the directory 'BaseDir'. Return error if anything goes wrong.
func createInstance(account Account, baseDir string) error {
	var instanceDirectory string = FolderPrefix + account.Name + account.Login

	if _, err := os.Stat(instanceDirectory); !os.IsNotExist(err) {
		fmt.Println(instanceDirectory, "already exist, skipping!")
		return nil
	}

	fmt.Println("creating instance ", instanceDirectory)

	if err := os.Mkdir(instanceDirectory, Permissions); err != nil {
		fmt.Errorf("error creating directory", instanceDirectory)
		return err
	}

	return recursiveCopy(baseDir, instanceDirectory)
}

func main() {
	Basedir := os.Getenv("BASE_DIR")
	if Basedir == "" {
		log.Fatal("BASE_DIR not set")
	}

	jsonFile, err := os.Open(ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var accounts []Account
	if err := json.Unmarshal(byteValue, &accounts); err != nil {
		log.Fatal(err)
	}

	for _, account := range accounts {
		err := createInstance(account, Basedir)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: work to launch each account MT5, generate config.ini for this
	}

	log.Println("successfully created MT5 instance for all accounts!")
}
