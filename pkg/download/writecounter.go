package download

import (
	"github.com/gosuri/uiprogress"
	"log"
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
		wt.UpdateProgress()
		wt.nextUpdate = (wt.Total/wt.ContentLength)*hundredthOfContentLen + hundredthOfContentLen
	}
	return n, nil
}

func (wt *WriteTracker) UpdateProgress() {
	if wt.ContentLength != 0 {
		err := wt.Pb.Set(wt.Total)
		if err != nil {
			log.Fatal(err)
		}
	}
}
