package layout

import "strings"

func Asset2017KeyboardText() string {
	return strings.Join([]string{
		"~1234567890_+",
		"`!@#$%^&*()-=",
		"QWDGJYPUL:{}|",
		"qwdgjypul;[]\\",
		"ASETFHNIOR\"",
		"asetfhnior'",
		"ZXCVBKM<>?",
		"zxcvbkm,./",
	}, "\n")
}