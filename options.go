package sendkeys

import (
	"time"
)

// KBOpt are KBWrap for our wrapper
type KBOpt func(*KBWrap)

// Stubborn will cause our sequences to continue despite errors.
// Otherwise, we will stop if our error count is over 0.
func Stubborn(o *KBWrap) {
	o.stubborn = true
}

// Noisy will cause all errors to be printed to stdout.
func Noisy(o *KBWrap) {
	o.noisy = true
}

// Random will use random sleeps throughout the typing process.
// Otherwise, a static 10 milliseconds will be used.
func Random(o *KBWrap) {
	o.random = true
}

// NoDelay will bypass the 2 second delay for linux, mostly for testing.
func NoDelay(o *KBWrap) {
	o.nodelay = true
}

// Delay allows to change the delay between keystrokes
func DelayBefore(delay time.Duration) KBOpt {
	return func(k *KBWrap) {
		k.beforeDuration = delay
	}
}

func KeystrokeDuration(d time.Duration) KBOpt {
	return func(k *KBWrap) {
		k.beforeDuration = d
	}
}

func DelayAfter(d time.Duration) KBOpt {
	return func(k *KBWrap) {
		k.afterDuration = d
	}
}

func WithKeyMap(keyMap KeyMap) KBOpt {
	return func(k *KBWrap) {
		k.keyMap = keyMap
	}
}
