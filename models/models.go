package models

import (
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteCounter struct {
	n   int // bytes read so far
	bar *pb.ProgressBar
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	wc.n += len(p)
	wc.bar.SetCurrent(int64(wc.n))
	return wc.n, nil
}

// Start progress bar
func (wc *WriteCounter) Start() {
	wc.bar.Start()
}

// Finish progress bar
func (wc *WriteCounter) Finish() {
	wc.bar.Finish()
}

// NewWriteCounter returns a new Write counter
func NewWriteCounter(total int) *WriteCounter {
	b := pb.New(total)
	b.SetRefreshRate(time.Second)
	b.Set(pb.Bytes, true)
	b.SetWriter(os.Stdout)

	return &WriteCounter{
		bar: b,
	}
}
