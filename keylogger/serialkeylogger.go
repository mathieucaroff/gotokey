package keylogger

import (
	"fmt"
	"time"

	"github.com/mathieucaroff/gotokey/layout"
	"github.com/moutend/go-hook/pkg/types"
)

func SerialKeyLogger(keyboardChan chan types.KeyboardEvent, keyboard layout.Keyboard) {
	var key types.KeyboardEvent

	var lastTime time.Time
	var now time.Time

	// var sign layout.Sign

	for {
		key = <-keyboardChan

		if key.Message == types.WM_KEYUP {
			continue
		}

		// sign = keyboard.Lower[key.ScanCode]

		now = time.Now()
		if lastTime.Unix()/60 != now.Unix()/60 {
			lastTime = now.Truncate(time.Minute)
			fmt.Printf("\n%s ", lastTime.Format(time.RFC3339)[:16])
		}

		fmt.Printf(" %d,", key.ScanCode)

		lastTime = now
	}
}
