package download

import (
	"fmt"
	"github.com/gosuri/uiprogress"
)

// WriteTracker counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteTracker struct {
	Pb            *uiprogress.Bar
	Total         int
	ContentLength int

	nextUpdate int
}

func (wt *WriteTracker) Write(p []byte) (int, error) {
	n := len(p)
	wt.Total += n

	if (wt.Total > wt.nextUpdate || wt.Total == wt.ContentLength) && wt.ContentLength != 0 {
		hundredthOfContentLen := wt.ContentLength / 100
		err := wt.UpdateProgress()
		if err != nil {
			return n, err
		}
		wt.nextUpdate = (wt.Total/wt.ContentLength)*hundredthOfContentLen + hundredthOfContentLen
	}
	return n, nil
}

func (wt *WriteTracker) UpdateProgress() error {
	if wt.ContentLength != 0 {
		err := wt.Pb.Set(wt.Total)
		if err != nil {
			return fmt.Errorf("error updating progress bar in writetracker: %v", err)
		}
	}
	return nil
}
