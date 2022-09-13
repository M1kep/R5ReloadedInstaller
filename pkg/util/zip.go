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

	pb := uiprogress.AddBar(len(archive.File)).AppendCompleted()
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

		if !strings.HasPrefix(filePath, filepath.Clean(destinationPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path '%s' while extracting '%s' from zip '%s'", filePath, f.Name, zipFile)
		}

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directory '%s' while extracting files from zip '%s': %v", filePath, zipFile, err)
			}
			pb.Incr()
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			if err != nil {
				return fmt.Errorf("error creating directory '%s' while extracting '%s' from zip '%s': %v", filepath.Dir(filePath), f.Name, zipFile, err)
			}
		}

		err = func() error {
			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("error opening destination file while extracting '%s' from zip '%s': %v", filePath, zipFile, err)
			}
			defer dstFile.Close()

			fileInArchive, err := f.Open()
			if err != nil {
				return fmt.Errorf("error opening file '%s' in zip '%s': %v", f.Name, zipFile, err)
			}
			defer fileInArchive.Close()

			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				return fmt.Errorf("error extracting file '%s' from zip '%s' to '%s': %v", f.Name, zipFile, dstFile.Name(), err)
			}
			pb.Incr()
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
