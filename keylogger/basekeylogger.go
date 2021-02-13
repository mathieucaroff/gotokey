package keylogger

import (
	"fmt"
	"strings"
	"time"

	"github.com/mathieucaroff/gotokey/keyutil"

	"github.com/moutend/go-hook/pkg/types"
)

func BaseKeyLogger(keyboardChan chan types.KeyboardEvent) {
	var key types.KeyboardEvent

	var lastTime time.Time
	var now time.Time
	var lastKey types.KeyboardEvent
	var silencedModifier types.VKCode

	text := ""
	textOutput := ""
	ok := true
	seconds := 0
	silence := false

	for {
		key = <-keyboardChan

		text, ok = VKCodeBaseNameMap[key.VKCode]
		if !ok {
			text = fmt.Sprintf("--%d/", key.VKCode)
		}
		if key.VKCode == silencedModifier && keyutil.IsDown(key) {
			if !silence {
				textOutput = text
			}
			text = ""
			silence = true
		} else if key.VKCode == lastKey.VKCode && keyutil.IsDown(lastKey) && keyutil.IsUp(key) && !silence {
			switch len(text) {
			case 1:
				textOutput = text + strings.ToLower(text)
			case 2:
				textOutput = fmt.Sprintf(":%s", text[1:])
			case 3:
				textOutput = fmt.Sprintf("=%s", text[1:])
			default:
				textOutput = fmt.Sprintf("==%s", text[2:])
			}
			text = ""
			silencedModifier = 0
			silence = false
		} else if key.Message == types.WM_KEYUP {
			switch len(text) {
			case 1:
				textOutput = strings.ToLower(text)
			case 2:
				textOutput = fmt.Sprintf("'%s", text[1:])
			case 3:
				textOutput = fmt.Sprintf("^%s", text[1:])
			default:
				textOutput = fmt.Sprintf("^^%s", text[2:])
			}
			text = ""
			silencedModifier = 0
			silence = false
		} else if keyutil.IsDown(key) && keyutil.IsModifierKey(key) {
			// silence nexts
			silencedModifier = key.VKCode
			silence = false
		} else {
			// clear silence
			silencedModifier = 0
			silence = false
		}
		lastKey = key

		if len(textOutput) > 0 {
			now = time.Now()
			if lastTime.Unix()/60 != now.Unix()/60 {
				lastTime = now.Truncate(time.Minute)
				fmt.Printf("\n%s ", lastTime.Format(time.RFC3339)[:16])
			}
			seconds = int(now.Sub(lastTime).Seconds())
			if seconds > 0 {
				fmt.Printf("%d", seconds)
				lastTime = now
			}

			fmt.Print(textOutput)
		}
		textOutput = text
	}
}
