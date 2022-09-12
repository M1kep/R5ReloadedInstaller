package download

import (
	"R5ReloadedInstaller/pkg/progress"
	"R5ReloadedInstaller/pkg/util"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Built off of https://golang.doc.xuwenliang.com/download-a-file-with-progress/

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(filepath string, url string, downloadMessage string) error {
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return fmt.Errorf("error creating tmp file for download: %v", err)
	}

	contentLength, err := util.GetContentLengthFromURL(url)
	if err != nil {
		return fmt.Errorf("error retrieving content length for file download: %v", err)
	}

	pb := progress.NewProgressBarWithMessage(downloadMessage, contentLength)
	// Get the data
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("GET request for %s while performing file download failed: %v", url, err)
	}

	// Create our progress reporter and pass it to be used alongside our writer
	writeTracker := &WriteTracker{
		Pb:            pb,
		ContentLength: contentLength,
	}
	_, err = io.Copy(out, io.TeeReader(resp.Body, writeTracker))
	if err != nil {
		return fmt.Errorf("error while downloading file from %s: %v", url, err)
	}

	// If the content-length was 0, the progress bar needs to be manually incremented to indicate completion
	if contentLength == 0 {
		pb.Incr()
	}
	out.Close()

	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}
