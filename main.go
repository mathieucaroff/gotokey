package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mathieucaroff/gotokey/keylogger"

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

	fmt.Println("ready")

	go keylogger.BaseKeyLogger(keyboardChan)

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
