package download

import (
	"github.com/gosuri/uiprogress"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteCounter struct {
	Pb            *uiprogress.Bar
	Total         int
	ContentLength int

	nextUpdate int
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += n

	if wc.Total > wc.nextUpdate && wc.ContentLength != 0 {
		wc.UpdateProgress()
		wc.nextUpdate = wc.nextUpdate + (wc.ContentLength / 100)
	}
	return n, nil
}

func (wc *WriteCounter) UpdateProgress() {
	if wc.ContentLength != 0 {
		wc.Pb.Incr()
		//err := wc.Pb.Set(wc.Total)
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
}
