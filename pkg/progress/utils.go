package progress

import (
	"github.com/gosuri/uiprogress"
)

// NewProgressBarWithMessage Will create a progress bar with the provided message
// Progress bars total will be set to 1 if maxProgress is 0
func NewProgressBarWithMessage(message string, maxProgress int) *uiprogress.Bar {
	var pb *uiprogress.Bar
	if maxProgress == 0 {
		pb = uiprogress.AddBar(1).AppendCompleted()
	} else {
		pb = uiprogress.AddBar(maxProgress).AppendCompleted()
	}
	pb.PrependFunc(func(b *uiprogress.Bar) string {
		return message
	})

	return pb
}
