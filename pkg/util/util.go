package util

import (
	"errors"
	"os"
)

func Exists(fileOrDirName string) (bool, error) {
	if _, err := os.Stat(fileOrDirName); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
