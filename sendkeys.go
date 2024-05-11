package sendkeys

import (
	"errors"
	"runtime"
	"sync"
	"time"

	kbd "github.com/micmonay/keybd_event"
)

// KBWrap is a wrapper for the keybd_event library for convenience
type KBWrap struct {
	d              kbd.KeyBonding
	errors         []error
	stubborn       bool
	noisy          bool
	random         bool
	nodelay        bool
	beforeDuration time.Duration
	downDuration   time.Duration // how long the key is pressed
	afterDuration  time.Duration // how long to wait after a key press

	mu sync.Mutex
}

func newKbw() *KBWrap {
	return &KBWrap{
		errors:         []error{},
		stubborn:       false,
		noisy:          false,
		random:         false,
		nodelay:        false,
		beforeDuration: 0 * time.Millisecond,
		downDuration:   40 * time.Millisecond,
		afterDuration:  10 * time.Millisecond,
	}
}

// NewKBWrapWithOptions creates a new keyboard wrapper with the given options.
// As of writing, those options include: Stubborn Noisy and Random.
// The defaults are all false.
func NewKBWrapWithOptions(opts ...KBOpt) (kbw *KBWrap, err error) {
	kbw = newKbw()
	kbw.d, err = kbd.NewKeyBonding()
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(kbw)
	}

	kbw.linDelay()
	return
}

func (kb *KBWrap) linDelay() {
	if kb.nodelay {
		return
	}
	// For linux, it is very important to wait 2 seconds
	// kayos note: idfk why tho, this is according to keybd_event author
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}
}

func (kb *KBWrap) down() {
	if !kb.check() {
		return
	}
	kb.handle(kb.d.Press())
}
func (kb *KBWrap) up() {
	if !kb.check() {
		return
	}
	kb.handle(kb.d.Release())
}

// press presses a key from the queue, waits, and then releases.
// Default wait time is 10 milliseconds.
func (kb *KBWrap) press() {
	kb.down()
	time.Sleep(kb.downDuration)
	kb.up()
}

func (kb *KBWrap) set(keys ...int) {
	kb.d.SetKeys(keys...)
}

func (kb *KBWrap) clr() {
	kb.d.Clear()
}

func (kb *KBWrap) only(k int) {
	kb.clr()
	kb.set(k)
	kb.press()
	kb.clr()
}

// Escape presses the escape key.
// All other keys will be cleared.
func (kb *KBWrap) Escape() {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	kb.only(kbd.VK_ESC)
}

// Tab presses the tab key.
// All other keys will be cleared.
func (kb *KBWrap) Tab() {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	kb.only(kbd.VK_TAB)
}

// Enter presses the enter key.
// All other keys will be cleared.
func (kb *KBWrap) Enter() {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	kb.only(kbd.VK_ENTER)
}

// BackSpace presses the backspace key.
// All other keys will be cleared.
func (kb *KBWrap) BackSpace() {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	kb.only(backspace)
}

// Type types out a string by simulating keystrokes.
// Check the exported Symbol map for non-alphanumeric keys.
func (kb *KBWrap) Type(s string) error {
	keys := kb.strToKeys(s)
	if !kb.check() {
		return errors.Join(kb.errors...)
	}

	kb.mu.Lock()
	defer kb.mu.Unlock()

	for _, key := range keys {
		kb.d.HasALT(key.ModifierALT)
		kb.d.HasSuper(key.ModifierSuper)
		kb.d.HasCTRL(key.ModifierCTRL)
		kb.d.HasSHIFT(key.ModifierSHIFT)
		kb.set(key.Code)
		kb.press()
		kb.clr()
	}
	return nil
}
