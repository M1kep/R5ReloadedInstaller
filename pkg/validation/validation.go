package validation

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsR5Folder(path string) (bool, error) {
	filesInDir, err := os.ReadDir(path)
	if err != nil {
		return false, fmt.Errorf("error while reading path '%s': %v", path, err)
	}

	for _, file := range filesInDir {
		if file.Name() == "r5apex.exe" {
			return true, nil
		}
	}

	return false, nil
}

func IsRunningInR5Folder() (bool, error) {
	path, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("error retrieving current directory while validating r5Path: %v", err)
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
