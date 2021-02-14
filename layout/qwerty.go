package layout

import "strings"

func QwertyKeyboardText() string {
	return strings.Join([]string{
		"~1234567890_+",
		"`!@#$%^&*()-=",
		"QWERTYUIOP{}|",
		"qwertyuiop[]\\",
		"ASDFGHJKL:\"",
		"asdfghjkl;'",
		"ZXCVBNM<>?",
		"zxcvbnm,./",
	}, "\n")
}
