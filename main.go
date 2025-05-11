package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
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
		fmt.Errorf("error opening file", err)
		return err
	}

	defer sourceFile.Close()

	destinationFile, err := os.Create(destination)
	if err != nil {
		fmt.Errorf("error creating file", err)
		return err
	}

	defer destinationFile.Close()

	if err := os.Chmod(destination, Permissions); err != nil {
		fmt.Errorf("error setting file permission", err)
		return err
	}

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

// createInstance attempts to create a directory. Its name prefixed with 'FolderPrefix'
// and some details in 'account'. The aim of this directory is to replicate everything
// fould in the directory 'BaseDir'. Return error if anything goes wrong.
func createInstance(account Account, baseDir string, instanceDirectory string) error {

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


// createInstanceConfig attempts to create the config.ini file start up file which will
// be used to boot up the MT5 instance in 'instanceDirectory' for 'account'
func createInstanceConfig(account Account, instanceDirectory string) (string, error) {
	var accountConfig string = instanceDirectory + "/" + account.Name + ".ini"
	configContent := fmt.Sprintf("[StartUp]\nLogin=%s\nPassword=%s\nServer=%s",
		account.Login, account.Password, account.Server)

	err := os.WriteFile(accountConfig, []byte(configContent), Permissions)
	if err != nil {
		fmt.Errorf("error writing config", accountConfig)
		return accountConfig, err
	}

	log.Println("config", accountConfig, "setup!")
	return accountConfig, err
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
		var instanceDirectory string = FolderPrefix + account.Name
		err := createInstance(account, Basedir, instanceDirectory)
		if err != nil {
			log.Fatal(err)
		}

		accountConfig, err := createInstanceConfig(account, instanceDirectory)
		if err != nil {
			log.Fatal(err)
		}

		// execute MT5 launcher
		command := exec.Command(account.Path, "/config:", accountConfig)
		if err := command.Start(); err != nil {
			log.Fatal("Failed to launch MT5 instance ", err)
		}

		log.Println(account.Path, "MT5 instance successfully started!")
	}

	log.Println("successfully created MT5 instance for all accounts!")
}
