package layout

import "strings"

func QwertyKeyboardText() string {
	return strings.Join([]string{
		"~1234567890_+",
		"`!@#$%^&*()-=",
		"QWERTYUIOP{}|",
		"qwertyuiop[]\\",
		"asdfghjkl:\"",
		"asdfghjkl;'",
		"ZXCVBNM<>?",
		"zxcvbnm,./",
	}, "\n")
}
