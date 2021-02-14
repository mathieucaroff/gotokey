package layout

import "strings"

func AzertyKeyboardText() string {
	return strings.Join([]string{
		" 1234567890°+",
		"²&é\"'(-è_çà)=",
		"AZERTYUIOP¨£µ",
		"azertyuiop^$*",
		"QSDFGHJKLM%",
		"qsdfghjklmù'",
		"WXCVBN?./§",
		"wxcvbn,;:!",
	}, "\n")
}
