package utils

import (
	"runtime"
	"time"

	"gopkg.in/cheggaaa/pb.v1"
)

func NewProgress(prefix string, count int) *pb.ProgressBar {
	// get the new original progress bar.
	bar := pb.New(count).Prefix(prefix)

	// Refresh rate for progress bar is set to 100 milliseconds.
	bar.SetRefreshRate(time.Millisecond * 100)

	// Use different unicodes for Linux, OS X and Windows.
	switch runtime.GOOS {
	case "linux":
		// Need to add '\x00' as delimiter for unicode characters.
		bar.Format("┃\x00▓\x00█\x00░\x00┃")
	case "darwin":
		// Need to add '\x00' as delimiter for unicode characters.
		bar.Format(" \x00▓\x00 \x00░\x00 ")
	default:
		// Default to non unicode characters.
		bar.Format("[=> ]")
	}
	bar.Start()
	return bar
}
