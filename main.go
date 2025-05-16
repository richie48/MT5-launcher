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
	Login    uint32
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

// createInstanceConfig attempts to create config file at 'expectedConfigLocation' which will be
// used as start up file when booting up the MT5 instance for 'account'. Returns an error if it
// fails to create config.
func createInstanceConfig(account Account, expectedConfigLocation string) error {
	configContent := fmt.Sprintf("[Login]\nLogin=%d\nPassword=%s\nServer=%s",
		account.Login, account.Password, account.Server)

	err := os.WriteFile(expectedConfigLocation, []byte(configContent), Permissions)
	if err != nil {
		fmt.Errorf("error writing config", expectedConfigLocation)
		return err
	}

	log.Println("config", expectedConfigLocation, "setup!")
	return err
}

// createInstance attempts to create a directory. Its name prefixed with 'FolderPrefix'
// and some details in 'account'. The aim of this directory is to replicate everything
// fould in the directory 'BaseDirectory'. Return error if anything goes wrong. Also returns
// name location of generated config ini
func createInstance(account Account, baseDirectory string, instanceDirectory string) (string, error) {

	var expectedConfigLocation string = filepath.Join(instanceDirectory, account.Name+".ini")

	if _, err := os.Stat(instanceDirectory); !os.IsNotExist(err) {
		fmt.Println(instanceDirectory, "already exist, skipping!")
		return expectedConfigLocation, nil
	}

	fmt.Println("creating instance ", instanceDirectory)

	if err := os.Mkdir(instanceDirectory, Permissions); err != nil {
		fmt.Errorf("error creating directory", instanceDirectory)
		return "", err
	}

	if err := createInstanceConfig(account, expectedConfigLocation); err != nil {
		fmt.Errorf("failed to generate config in location", instanceDirectory)
		return expectedConfigLocation, err
	}

	return expectedConfigLocation, recursiveCopy(baseDirectory, instanceDirectory)
}

func main() {
	BaseDirectory := os.Getenv("BASE_DIR")
	if BaseDirectory == "" {
		log.Fatal("BASE_DIR not set")
	}

	SourceDirectory := os.Getenv("SRC_DIR")
	if SourceDirectory == "" {
		log.Fatal("SRC_DIR not set")
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

	noLaunch := false // default
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--no-launch" {
			noLaunch = true
			log.Println("running in NO_LAUNCH MODE")
		}
	}

	for _, account := range accounts {
		var fullSourceDirectoryPath string = filepath.Join(BaseDirectory, SourceDirectory)
		var fullInstanceDirectoryPath string = filepath.Join(BaseDirectory, FolderPrefix+account.Name)
		accountConfig, err := createInstance(account, fullSourceDirectoryPath, fullInstanceDirectoryPath)
		if err != nil {
			log.Fatal(err)
		}

		if !noLaunch {
			configArgument := "/config:" + accountConfig
			command := exec.Command(account.Path, configArgument)
			if err := command.Start(); err != nil {
				log.Fatal("Failed to launch MT5 instance ", err)
			}
			log.Println(account.Path, "MT5 instance successfully started!")
		}

	}

	log.Println("successfully created MT5 instance for all accounts!")
}
