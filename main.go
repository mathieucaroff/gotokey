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

	fmt.Println("start capturing keyboard input")

	var lastTime time.Time
	var now time.Time
	var lastKey types.KeyboardEvent

	text := ""
	textOutput := ""
	ok := true
	seconds := 0

	for {
		select {
		case <-time.After(5 * time.Minute):
			fmt.Println("Received timeout signal")
			return nil
		case <-signalChan:
			fmt.Println("Received shutdown signal")
			return nil
		case key := <-keyboardChan:
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

			text, ok = VKCodeMap[key.VKCode]
			if !ok {
				text = fmt.Sprintf("--0x%X/", key.VKCode)
			}
			if key.VKCode == lastKey.VKCode &&
				lastKey.Message == types.WM_KEYDOWN &&
				key.Message == types.WM_KEYUP {
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
			}
			lastKey = key

			fmt.Print(textOutput)
			textOutput = text
		}
	}
}
