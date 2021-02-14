package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/mathieucaroff/gotokey/keylogger"
	"github.com/mathieucaroff/gotokey/layout"

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
	keyboardChannel := make(chan types.KeyboardEvent, 100)

	if err := keyboard.Install(nil, keyboardChannel); err != nil {
		return err
	}

	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("ready")

	loggerRoutine := keylogger.RawKeyLogger
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "raw":
			loggerRoutine = keylogger.RawKeyLogger
		case "base":
			loggerRoutine = keylogger.BaseKeyLogger
		case "serial":
			keyboardText := ""
			if len(os.Args) > 2 {
				kbArg := os.Args[2]
				if strings.Count(kbArg, "\n") == 7 {
					keyboardText = kbArg
				} else if kbArg == "asset" {
					keyboardText = layout.Asset2017KeyboardText()
				} else {
					fmt.Fprintf(os.Stderr, "unrecognized %s\n", kbArg)
					os.Exit(1)
				}
			}
			if len(keyboardText) == 0 {
				keyboardText = layout.QwertyKeyboardText()
			}
			keyboard := layout.KeyboardFromText(keyboardText)
			loggerRoutine = func(c chan types.KeyboardEvent) {
				keylogger.SerialKeyLogger(c, keyboard)
			}
		}
	}

	go loggerRoutine(keyboardChannel)

	for {
		select {
		case <-time.After(5 * time.Minute):
			fmt.Println("\nTimeout")
		case <-signalChan:
			fmt.Println("\nShutdown")
		}
		os.Exit(0)
	}
}
