package sendkeys

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/eiannone/keyboard"
)

func Test_strToKeys(t *testing.T) {
	strTo("yeet", t)
	strTo("YEET", t)
	strTo("YeeT", t)
}

func Test_NewKBWrapWithOptions(t *testing.T) {
	k, err := NewKBWrapWithOptions(Noisy, NoDelay, Stubborn, Random)
	if err != nil {
		t.Fatal(err.Error())
	}
	opts := []*bool{&k.noisy, &k.nodelay, &k.stubborn, &k.random}
	for _, opt := range opts {
		if *opt != true {
			t.Fatalf("KBWrap should have had options Noisy: true, NoDelay: true, Stubborn: true, Random: true. "+
				"Had Noisy: %t NoDelay: %t Stubborn: %t Random: %t", k.noisy, k.nodelay, k.stubborn, k.random)
		}
	}
	t.Logf("[OPT] Noisy: %t NoDelay: %t Stubborn: %t Random: %t", k.noisy, k.nodelay, k.stubborn, k.random)
	k = nil
	opts = nil
	k, err = NewKBWrapWithOptions()
	if err != nil {
		t.Fatal(err.Error())
	}
	opts = []*bool{&k.noisy, &k.nodelay, &k.stubborn, &k.random}
	for _, opt := range opts {
		if *opt != false {
			t.Fatalf("KBWrap should have had options Noisy: false, NoDelay: false, Stubborn: false, Random: false. "+
				"Had Noisy: %t NoDelay: %t Stubborn: %t Random: %t", k.noisy, k.nodelay, k.stubborn, k.random)
		}
	}
	t.Logf("[OPT] Noisy: %t NoDelay: %t Stubborn: %t Random: %t", k.noisy, k.nodelay, k.stubborn, k.random)
	k = nil
	opts = nil
}

func TestSendkeysAlphaUpper(t *testing.T) {
	ret := make(chan rune)
	//go listenForKeys(t, ret)
	k, err := NewKBWrapWithOptions(Noisy)
	if err != nil {
		t.Fatal(err)
	}
	testsend(t, k, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", ret)
}

func TestSendkeysAlphaLower(t *testing.T) {
	ret := make(chan rune)
	//go listenForKeys(t, ret)
	k, err := NewKBWrapWithOptions(Noisy)
	if err != nil {
		t.Fatal(err)
	}
	testsend(t, k, "abcdefghijklmnopqrstuvwxyz", ret)
}

func TestSendkeysNumeric(t *testing.T) {
	ret := make(chan rune)
	//go listenForKeys(t, ret)
	k, err := NewKBWrapWithOptions(Noisy)
	if err != nil {
		t.Fatal(err)
	}
	testsend(t, k, "1234567890", ret)
}

func TestSendkeysSymbols(t *testing.T) {
	ret := make(chan rune)
	//go listenForKeys(t, ret)
	k, err := NewKBWrapWithOptions()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	testsend(t, k, "-_=+[{]}'\"`~\\|,<.>/? !@#$%^&*();:", ret)
}

func testsend(t *testing.T, k *KBWrap, teststr string, ret chan rune) {

	var (
		expected      = len([]rune(teststr))
		count         = 0
		chars         []rune
		ctx, cancel   = context.WithCancelCause(context.Background())
		cancelTimeout context.CancelFunc
	)
	ctx, cancelTimeout = context.WithTimeout(ctx, 5*time.Second)
	defer cancel(context.Canceled)
	defer cancelTimeout()

	go func(s string) {
		err := k.Type(s)
		if err != nil {
			cancel(err)
		}
	}(teststr)

loop:
	for {
		select {
		case <-ctx.Done():
			t.Error(ctx.Err())
			return
		case chr, ok := <-ret:
			if !ok {
				cancel(errors.New("ret channel closed unexpectedly"))
				t.Error(ctx.Err())
				return
			}
			chars = append(chars, chr)
			count++

			if len(chars) >= expected {
				break loop
			}
		}
	}

	var final = string(chars)
	if final != teststr {
		t.Errorf("[FAIL] Have: %s, Wanted: %s", final, teststr)
	} else {
		t.Logf(
			"[SUCCESS] got %d characters: %s",
			count, final,
		)
	}
}

func listenForKeys(t *testing.T, ret chan rune) {

	keysEvents, err := keyboard.GetKeys(1)
	if err != nil {
		t.Logf("failed to listen to keyboard: %v, closing ret channel", err)
		close(ret)
	}
	defer keyboard.Close()

	for event := range keysEvents {
		if event.Err != nil {
			t.Logf("failed to listen to keyboard events: %v, closing ret channel", event.Err)
			close(ret)
		}
		t.Logf("Key pressed: %v (%d)", event.Rune, event.Key)
		ret <- event.Rune
	}

}

func strTo(teststr string, t *testing.T) {
	split := strings.Split(teststr, "")
	k, err := NewKBWrapWithOptions(Noisy, NoDelay)
	if err != nil {
		t.Fatalf(err.Error())
	}
	keys := k.strToKeys(teststr)
	if len(keys) != len(split) {
		t.Fatalf("length of mapped keys: %d, wanted length of string: %d", len(keys), len(split))
	}
	t.Logf("string: %s, keys: %#v", teststr, keys)
}
