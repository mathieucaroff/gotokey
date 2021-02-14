package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

	fmt.Println("ready")

	var loggerRoutine func(c chan types.KeyboardEvent)

	args := os.Args
	if len(args) < 2 {
		args = append(args, "serial", "asset")
	}

	switch args[1] {
	case "raw":
		loggerRoutine = keylogger.RawKeyLogger
	case "base":
		loggerRoutine = keylogger.BaseKeyLogger
	case "serial":
		keyboardText := ""
		if len(args) > 2 {
			kbArg := args[2]
			if strings.Count(kbArg, "\n") == 7 {
				keyboardText = kbArg
			} else if kbArg == "asset" {
				keyboardText = layout.Asset2017KeyboardText()
			} else if kbArg == "azerty" {
				keyboardText = layout.AzertyKeyboardText()
			} else if kbArg == "qwerty" {
				keyboardText = layout.QwertyKeyboardText()
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

	loggerRoutine(keyboardChannel)

	return nil
}
