package download

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"io"
	"net/http"
	"os"
	"strconv"
)

// https://golang.doc.xuwenliang.com/download-a-file-with-progress/

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(filepath string, url string, downloadMessage string) error {
	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	headResp, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("HEAD request for %s failed: %v", url, err)
	}

	// If we are not provided a content lenght, the progress bar will be a 0/1 rather than actual progress.
	cntlenheader := headResp.Header.Get("Content-Length")
	cntlen := 0
	if cntlenheader != "" {
		cntlen, err = strconv.Atoi(cntlenheader)
		if err != nil {
			return fmt.Errorf("failed to convert \"%s\" to int: %v", cntlenheader, err)
		}
	}

	var pb *uiprogress.Bar
	if cntlen == 0 {
		pb = uiprogress.AddBar(1).AppendCompleted()
	} else {
		//pb = uiprogress.AddBar(cntlen).AppendCompleted()
		pb = uiprogress.AddBar(100).AppendCompleted()
	}
	pb.PrependFunc(func(b *uiprogress.Bar) string {
		return downloadMessage
	})

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("GET request for %s failed: %v", url, err)
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{
		Pb:            pb,
		ContentLength: cntlen,
	}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// Because we don't know the content length, this is just a started/completed tracker
	if cntlen == 0 {
		pb.Incr()
	}
	out.Close()

	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}
