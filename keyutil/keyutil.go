package keyutil

import (
	"github.com/moutend/go-hook/pkg/types"
)

func IsDown(key types.KeyboardEvent) bool {
	return key.Message == types.WM_KEYDOWN
}

func IsUp(key types.KeyboardEvent) bool {
	return key.Message == types.WM_KEYUP
}

func IsModifierKey(key types.KeyboardEvent) bool {
	return key.VKCode == types.VK_CONTROL || // Control
		key.VKCode == types.VK_LCONTROL ||
		key.VKCode == types.VK_RCONTROL ||
		key.VKCode == types.VK_SHIFT || // Shift
		key.VKCode == types.VK_LSHIFT ||
		key.VKCode == types.VK_RSHIFT ||
		key.VKCode == types.VK_MENU || // Alt
		key.VKCode == types.VK_LMENU ||
		key.VKCode == types.VK_RMENU ||
		key.VKCode == types.VK_LWIN || // Win
		key.VKCode == types.VK_RWIN ||
		false
}
