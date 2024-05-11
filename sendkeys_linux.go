package sendkeys

//go: build +linux
import kbd "github.com/micmonay/keybd_event"

const (
	backspace = kbd.VK_BACKSPACE
)
