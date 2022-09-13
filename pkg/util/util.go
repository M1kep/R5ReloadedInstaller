package util

import (
	"errors"
	"fmt"
	"github.com/tawesoft/golib/v2/dialog"
	"net/http"
	"os"
	"strconv"
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

func GetContentLengthFromURL(url string) (contentLength int, err error) {
	headResp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("HEAD request for %s failed: %v", url, err)
	}

	contentLengthHeader := headResp.Header.Get("Content-Length")
	if contentLengthHeader == "" {
		return 0, nil
	}

	contentLength, err = strconv.Atoi(contentLengthHeader)
	if err != nil {
		return 0, fmt.Errorf("failed to convert \"%s\" to int: %v", contentLengthHeader, err)
	}
	return contentLength, nil
}

func LogErrorWithDialog(err error) {
	fmt.Println(err)
	_ = dialog.Error("Program encountered error. See console for logs.")
}
