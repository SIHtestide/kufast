package trackerFactory

import (
	"github.com/jedib0t/go-pretty/v6/progress"
	"time"
)

func CreateProgressWriter(numExpected int) progress.Writer {
	pw := progress.NewWriter()
	pw.SetNumTrackersExpected(numExpected)
	pw.SetMessageWidth(100)
	pw.SetUpdateFrequency(time.Millisecond * 250)

	return pw
}

func HandleTracking(pw progress.Writer, expectedTracker int) {
	go pw.Render()

	for !pw.IsRenderInProgress() {
		time.Sleep(time.Millisecond * 200)
	}
	for pw.IsRenderInProgress() {
		time.Sleep(time.Millisecond * 100)
		if pw.LengthDone() == expectedTracker {
			pw.Stop()
		}
	}
}
