// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Buffer size is depends on your need. The 100 is placeholder value.
	keyboardChan := make(chan types.KeyboardEvent, 100)

	if err := keyboard.Install(nil, keyboardChan); err != nil {
		return err
	}

	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	var lastTime time.Time
	var now time.Time
	var lastKey types.KeyboardEvent
	var silencedModifier types.VKCode

	text := ""
	textOutput := ""
	ok := true
	seconds := 0
	silence := false

	fmt.Printf("%#x %d ", uint32(4), uint32(5))
	fmt.Println("ready")

	for {
		select {
		case <-time.After(5 * time.Minute):
			fmt.Println("Timeout")
			return nil
		case <-signalChan:
			fmt.Println("Shutdown")
			return nil
		case key := <-keyboardChan:
			text, ok = VKCodeMap[key.VKCode]
			if !ok {
				text = fmt.Sprintf("--%d/", key.VKCode)
			}
			if key.VKCode == silencedModifier && isDown(key) {
				if !silence {
					textOutput = text
				}
				text = ""
				silence = true
			} else if key.VKCode == lastKey.VKCode && isDown(lastKey) && isUp(key) && !silence {
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
			} else if isDown(key) && isModifierKey(key) {
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
}

func isDown(key types.KeyboardEvent) bool {
	return key.Message == types.WM_KEYDOWN
}

func isUp(key types.KeyboardEvent) bool {
	return key.Message == types.WM_KEYUP
}

func isModifierKey(key types.KeyboardEvent) bool {
	return key.VKCode == types.VK_CONTROL || // Control
		key.VKCode == types.VK_LCONTROL ||
		key.VKCode == types.VK_RCONTROL ||
		key.VKCode == types.VK_SHIFT || // Shift
		key.VKCode == types.VK_LSHIFT ||
		key.VKCode == types.VK_RSHIFT ||
		key.VKCode == types.VK_MENU || // Alt
		key.VKCode == types.VK_LMENU ||
		key.VKCode == types.VK_RMENU ||
		key.VKCode == types.VK_LWIN || // Win
		key.VKCode == types.VK_RWIN ||
		false
}
