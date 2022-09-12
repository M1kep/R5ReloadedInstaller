package util

import (
	"archive/zip"
	"fmt"
	"github.com/gosuri/uiprogress"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnzipFile(zipFile string, destinationPath string, stripFirstFolder bool, progressMessage string) error {
	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(fmt.Errorf("error opening zipfile %s: %v", zipFile, err))
	}
	defer archive.Close()

	pb := uiprogress.AddBar(len(archive.File))
	pb.PrependFunc(func(b *uiprogress.Bar) string {
		return progressMessage
	})

	for _, f := range archive.File {
		var filePath string
		if stripFirstFolder {
			fNameRemovedFirstPath := strings.Split(f.Name, "/")[1:]
			if fNameRemovedFirstPath[0] == "" {
				pb.Incr()
				continue
			}
			filePath = filepath.Join(destinationPath, strings.Join(fNameRemovedFirstPath, string(os.PathSeparator)))
		} else {
			filePath = filepath.Join(destinationPath, f.Name)
		}
		pb.Incr()

		if !strings.HasPrefix(filePath, filepath.Clean(destinationPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path")
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}
