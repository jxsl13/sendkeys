package sendkeys

import (
	"strings"
	"testing"
	"time"

	"github.com/eiannone/keyboard"
)

func listenForKeys(t *testing.T, ret chan string) {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		_ = keyboard.Close()
	})

	for {
		event := <-keysEvents
		if event.Err != nil {
			t.Fatalf(event.Err.Error())
		}
		t.Logf("Key pressed: %s", string(event.Rune))
		if event.Key == keyboard.KeyEsc {
			break
		}
		ret <- string(event.Rune)
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

func Test_sendkeys(t *testing.T) {
	ret := make(chan string)
	go listenForKeys(t, ret)
	k, err := NewKBWrapWithOptions(Noisy)
	if err != nil {
		t.Fatal(err)
	}
	testsend(t, k, "yeet", ret)
	testsend(t, k, "YEET", ret)
	testsend(t, k, "YeeT", ret)
	// testsend(t, k,"Yeet!")
}

func testsend(t *testing.T, k *KBWrap, teststr string, ret chan string) {
	keys := k.strToKeys(teststr)

	k.set(keys...)

	var count = 0
	var chars []string

	go func() {
		// t.Logf("[receiver(%s)] go func() start", teststr)
		// defer t.Logf("[receiver(%s)] go func() return", teststr)
		for {
			select {
			case chr := <-ret:
				chars = append(chars, chr)
				count++
			default:
				if count == len(teststr) {
					return
				}
			}
		}
	}()

	//	t.Log("sleeping for 250ms...")
	//	time.Sleep(250 * time.Millisecond)

	err := k.Type(teststr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	var tick = 0
	var brk = false
	for {
		if brk {
			break
		}
		switch {
		case tick >= 5:
			t.Logf("FAIL: took too long")
			t.Fail()
			brk = true
			break
		case count == len(teststr):
			brk = true
			break
		default:
			time.Sleep(1 * time.Second)
			tick++
		}
	}

	t.Logf(
		"got %d characters, got %s string.",
		count, strings.Join(chars, ""),
	)

	if strings.Join(chars, "") != teststr {
		t.Fail()
	}

}
