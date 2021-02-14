package layout

import (
	"fmt"
	"strings"
)

// 1, 59, 60, 61, 62, 63, 64, 65, 66, 67,68, 87, 88,
// 55, 83,
// 41, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
// 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 43,
// 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 28,
// 42, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54,
// 29, 91, 56, 57, 86, 541, 56, 29,
//
// 72, 75, 80, 77,
// 71, 73, 81, 79,

// A Keyboard maps each keycode for uppercase and lowercase inserts to a
// Sign
type Keyboard struct {
	Upper []Sign
	Lower []Sign
}

// SignKind -
type SignKind uint8

const (
	NoSignKind SignKind = iota
	Alpha
	Numeric
	Punctuation
	Move
	Space
	Tabulation
	Enter
	Backspace
	Delete
	PrintScreen
	ScrollLock
	Pause
	CapsLock
	Escape
	FKey
	Modifier
)

// A Sign describes a keyboard input that can be sent to an application
type Sign struct {
	Name string
	Kind SignKind
}

func MakeSign(a byte) Sign {
	var kind SignKind
	if ('A' <= a && a <= 'Z') || ('a' <= a && a <= 'z') {
		kind = Alpha
	} else if '0' <= a && a <= '9' {
		kind = Numeric
	} else {
		kind = Punctuation
	}
	return Sign{
		Name: string(a),
		Kind: kind,
	}
}

// KeyboardFromText This copies the characters in a text description of a keyboard into an array
// mapping scancode to value
func KeyboardFromText(text string) Keyboard {
	grid := strings.Split(text, "\n")

	upper := make([]Sign, 128)
	lower := make([]Sign, 128)

	Register := func(kind SignKind) func(k int, name string) {
		return func(k int, name string) {
			upper[k] = Sign{Kind: kind, Name: strings.ToUpper(name)}
			lower[k] = Sign{Kind: kind, Name: strings.ToLower(name)}
		}
	}

	// E row, [1:16]
	for k := 1; k < 14; k++ {
		upper[k] = MakeSign(grid[0][k-1])
		lower[k] = MakeSign(grid[1][k-1])
	}
	Register(Backspace)(14, "backspace")

	// D row, [16:30]
	Register(Tabulation)(15, "tabulation")
	for k := 16; k < 28; k++ {
		upper[k] = MakeSign(grid[2][k-16])
		lower[k] = MakeSign(grid[3][k-16])
	}
	upper[43] = MakeSign(grid[2][28-16])
	lower[43] = MakeSign(grid[3][28-16])

	// C row, [29:40]+[28]
	Register(CapsLock)(29, "capslock")
	for k := 30; k < 41; k++ {
		upper[k] = MakeSign(grid[4][k-30])
		lower[k] = MakeSign(grid[5][k-30])
	}
	Register(Enter)(28, "enter")

	// B row, [42:55]
	Register(Modifier)(42, "lshift")
	for k := 44; k < 54; k++ {
		upper[k] = MakeSign(grid[6][k-44])
		lower[k] = MakeSign(grid[7][k-44])
	}
	Register(Modifier)(54, "rshift")

	// A row, [29, 91, 56, 57, 86, (541, 56)]
	Register(Modifier)(29, "lcontrol")
	Register(Modifier)(91, "lwin")
	Register(Modifier)(56, "lalt")
	Register(Space)(57, "space")
	upper[86] = MakeSign('|')
	lower[86] = MakeSign('\\')
	Register(Modifier)(541, "ralt")

	// G1 block, [64:82]
	Register(Move)(75, "left")
	Register(Move)(77, "right")
	Register(Move)(72, "up")
	Register(Move)(80, "down")

	Register(Move)(71, "home")
	Register(Move)(79, "end")
	Register(Move)(73, "pgup")
	Register(Move)(81, "pgdown")

	Register(PrintScreen)(55, "printscreen")
	Register(ScrollLock)(83, "scrolllock")

	// G2 block, [75:90]
	// Register(Move)(79, "left")
	// Register(Move)(89, "right")
	// Register(Move)(83, "up")
	// Register(Move)(84, "down")

	// Register(Move)(80, "home")
	// Register(Move)(81, "end")
	// Register(Move)(85, "pgup")
	// Register(Move)(86, "pgdown")
	// Register(PrintScreen)(75, "printscreen")
	// Register(ScrollLock)(76, "scrolllock")
	// Register(Delete)(77, "delete")

	// H block, [90:109]

	// TODO: H block (numpad)

	// F row, [110:126]
	Register(Escape)(110, "escape")
	for k := 112; k < 126; k++ {
		Register(FKey)(k, fmt.Sprintf("f%d", k-112))
	}

	return Keyboard{
		Upper: upper,
		Lower: lower,
	}
}
