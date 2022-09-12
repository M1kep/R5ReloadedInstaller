package validation

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func IsR5Folder(path string) bool {
	filesInDir, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range filesInDir {
		if file.Name() == "r5apex.exe" {
			return true
		}
	}

	return false
}

func IsInR5Folder() bool {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(fmt.Errorf("error retrieving current directory while validating r5Path: %v", err))
	}

	return IsR5Folder(path)
}

func IsLauncherFileLocked(path string) bool {
	launcherPath := filepath.Join(path, "launcher.exe")
	file, err := os.OpenFile(launcherPath, os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return true
	}

	return false
}
