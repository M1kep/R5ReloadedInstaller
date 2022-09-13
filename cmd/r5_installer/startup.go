package main

import (
	"R5ReloadedInstaller/pkg/validation"
	"fmt"
	"github.com/tawesoft/golib/v2/dialog"
	"os"
	"path/filepath"
)

func getValidatedR5Folder() (validatedFolder string, err error) {
	isRunningInR5Folder, err := validation.IsRunningInR5Folder()
	if err != nil {
		return "", err
	}

	if isRunningInR5Folder {
		validatedFolder, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("error retrieving current directory while validating r5Path: %v", err)
		}

		return validatedFolder, nil
	}

	// Check CLI argument
	if !(len(os.Args) >= 2) {
		_ = dialog.Raise("Please move the R5RInstaller into your R5 Directory")
		return "", fmt.Errorf("not running from r5Folder and no argument provided")
	}

	pathFromArgs := os.Args[1]
	isPathR5Folder, err := validation.IsR5Folder(pathFromArgs)
	if err != nil {
		return "", err
	}
	if isPathR5Folder {
		validatedFolder = pathFromArgs
	} else {
		_ = dialog.Raise("Please move the R5RInstaller into your R5 Directory or pass correct path via arguments")
		return "", fmt.Errorf("not running from r5Folder and provided path argument is invalid")
	}

	return validatedFolder, nil
}

func initializeDirectories(r5Folder string) (cacheDir string, err error) {
	cacheDir = filepath.Join(r5Folder, "R5InstallerDirectory/cache")

	err = os.MkdirAll(cacheDir, 0777)
	if err != nil {
		err = fmt.Errorf("error initializing installer directory %s: %v", cacheDir, err)
		return
	}

	return
}
