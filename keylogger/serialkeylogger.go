package keylogger

import (
	"fmt"
	"strings"
	"time"

	"github.com/mathieucaroff/gotokey/layout"
	"github.com/moutend/go-hook/pkg/types"
)

type KeyState bool

const (
	UP   KeyState = false
	DOWN KeyState = true
)

func SerialKeyLogger(keyboardChan chan types.KeyboardEvent, keyboard layout.Keyboard) {
	var key types.KeyboardEvent

	var lastTime time.Time
	var now time.Time

	lcontrol := UP
	rcontrol := UP
	lwin := UP
	rwin := UP
	lalt := UP
	ralt := UP
	lshift := UP
	rshift := UP

	control := false
	win := false
	alt := false
	shift := false

	modifierString := ""

	setModifier := func(name string, state KeyState) {
		// set the named modifier to the given state
		switch name {
		case "lcontrol":
			lcontrol = state
		case "rcontrol":
			rcontrol = state
		case "lwin":
			lwin = state
		case "rwin":
			rwin = state
		case "lalt":
			lalt = state
		case "ralt":
			ralt = state
		case "lshift":
			lshift = state
		case "rshift":
			rshift = state
		default:
			return
		}

		control = bool(lcontrol) || bool(rcontrol)
		win = bool(lwin) || bool(rwin)
		alt = bool(lalt) || bool(ralt)
		shift = bool(lshift) || bool(rshift)

		// recompute the modifier string
		modifierString = ""
		if control {
			modifierString += "C"
		}
		if win {
			modifierString += "W"
		}
		if alt {
			modifierString += "A"
		}
	}

	var sign layout.Sign

	stdout := make(chan string)
	formatted := make(chan string)
	content := make(chan string)
	namedContent := make(chan string)

	// Print stuff (stdout -> os.stdout)
	go func() {
		for {
			fmt.Print(<-stdout)
		}
	}()

	// Timestamp stuff before printing it (formatted -> stdout)
	go func() {
		piece := ""
		for {
			piece = <-formatted
			now = time.Now()
			if lastTime.Unix()/60 != now.Unix()/60 {
				lastTime = now.Truncate(time.Minute)
				stdout <- fmt.Sprintf("\n%s|", lastTime.Format(time.RFC3339)[:16])
			}
			stdout <- piece
			lastTime = now
		}
	}()

	// Prefix pieces with modifier (content -> formatted)
	go func() {
		piece := ""
		for {
			piece = <-content
			if len(modifierString) > 0 {
				if piece[0] == '.' {
					formatted <- fmt.Sprintf(".%s%s", modifierString, piece)
				} else {
					formatted <- fmt.Sprintf(".%s.%s/", modifierString, piece)
				}
			} else {
				formatted <- piece
			}
		}
	}()

	// Add decorations around named keys (namedContent -> content)
	go func() {
		piece := ""
		for {
			piece = <-namedContent
			if shift {
				piece = strings.ToUpper(piece)
			}
			content <- fmt.Sprintf(".%s/", piece)
		}
	}()

	// Listen to user
	for {
		key = <-keyboardChan

		// Ignore
		if key.Message == types.WM_KEYUP {
			name := keyboard.Lower[key.ScanCode].Name
			setModifier(name, UP)
			continue
		}

		if shift {
			sign = keyboard.Upper[key.ScanCode]
		} else {
			sign = keyboard.Lower[key.ScanCode]
		}

		if sign.Kind == layout.Alpha || sign.Kind == layout.Numeric {
			content <- sign.Name
		} else if sign.Kind == layout.Punctuation {
			if sign.Name == "." {
				namedContent <- "dot"
			} else if sign.Name == "|" {
				namedContent <- "pipe"
			} else {
				content <- sign.Name
			}
		} else if sign.Kind == layout.Modifier {
			name := keyboard.Lower[key.ScanCode].Name
			setModifier(name, DOWN)
		} else if sign.Kind == layout.Space {
			content <- " "
		} else if len(sign.Name) > 0 {
			content <- fmt.Sprintf(".%s/", sign.Name)
		}
	}
}
