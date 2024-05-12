package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jxsl13/sendkeys"
	"github.com/manifoldco/promptui"
)

func main() {
	resultMap := map[string]sendkeys.KeyCode{}
	defer func() {
		if i := recover(); i != nil {
			fmt.Println(i)
		}

		data, _ := json.MarshalIndent(resultMap, "", " ")
		fmt.Println(string(data))
	}()

	for promptContinue("Continue filling the key codes map?") {
		cfg, err := configure()
		if err != nil {
			log.Println(err)
			return
		}

		code, err := selectMultiStep(cfg)
		if err != nil {
			log.Println(err)
			return
		}

		r, err := promptRune()
		if err != nil {
			log.Println(err)
			return
		}
		resultMap[string(r)] = code
	}
}

func selectMultiStep(cfg Config) (sendkeys.KeyCode, error) {
	codes := sendkeys.GenerateKeyCodesWithTemplate(
		cfg.Template,
		cfg.Start,
		cfg.End+1,
	)

	kb, err := sendkeys.NewKBWrapWithOptions(
		sendkeys.KeystrokeDuration(15*time.Millisecond),
		sendkeys.DelayAfter(50*time.Millisecond),
	)
	if err != nil {
		return sendkeys.KeyCode{}, fmt.Errorf("failed to initialize sender: %w", err)
	}

	for len(codes) > 1 {
		if !promptContinue("Start next typing step?") {
			return sendkeys.KeyCode{}, errors.New("aborted by user")
		}

		codes, err = selectSingleStep(kb, codes)
		if err != nil {
			return sendkeys.KeyCode{}, fmt.Errorf("failed to do single step: %w", err)
		}
	}

	if len(codes) == 0 {
		return sendkeys.KeyCode{}, errors.New("key code not found, no codes left in list")
	}

	code := codes[0]
	if promptContinue(fmt.Sprintf("Print final character one last time (%s)?", code)) {
		kb.TypeRaw(code)
	}

	return codes[0], nil

}

func selectSingleStep(kb *sendkeys.KBWrap, codes []sendkeys.KeyCode) ([]sendkeys.KeyCode, error) {
	var (
		start = 0
		end   = len(codes)
		half  = start + (end-start)/2
	)

	fmt.Printf("start: %s\n  end: %s\n", codes[start], codes[end-1])
	fmt.Println("starting step in 5 seconds...")
	time.Sleep(5 * time.Second)

	for i := start; i < end; i++ {
		if i == half {
			fmt.Println("First half printed")
			time.Sleep(5 * time.Second)
		}
		code := codes[i]
		kb.TypeRaw(code)
	}

	h, err := selectHalf()
	if err != nil {
		return nil, err
	}

	if h == 0 {
		return codes[start:half], nil
	}
	return codes[half:end], nil

}

func selectHalf() (int, error) {
	selected := 0
	prompt := promptui.Prompt{
		Label: "Which half contains your looked for character or symbol, 1 or 2?",
		Validate: func(input string) error {
			i, err := strconv.ParseInt(input, 10, 8)
			if err != nil {
				return errors.New("invalid number, must be 1 or 2")
			}

			switch i {
			case 1, 2:
			default:
				return errors.New("invalid number, must be 1 or 2")
			}
			selected = int(i)
			return nil
		},
	}

	_, err := prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("failed to prompt user: %w", err)
	}
	return selected - 1, nil
}

type Config struct {
	Start    int
	End      int
	Template sendkeys.KeyCode
}

func configure() (Config, error) {
	var (
		err error
		// underlying library is buggy when all
		// prompts are confirmed quickly
		// that's why we synchronize the validation below
		mu sync.Mutex
	)

	cfg := Config{
		Start:    0,
		End:      96,
		Template: sendkeys.SimpleKeyCode(0),
	}
	prompt := promptui.Prompt{
		Default:   strconv.FormatInt(int64(cfg.Start), 10),
		AllowEdit: true,
		Label:     "Start key code",
		Validate: func(input string) error {
			mu.Lock()
			defer mu.Unlock()

			i, err := strconv.ParseInt(input, 10, 16)
			if err != nil {
				return err
			}
			cfg.Start = int(i)
			return nil
		},
	}

	_, err = prompt.Run()
	if err != nil {
		return cfg, fmt.Errorf("failed to prompt start key code: %w", err)
	}

	prompt = promptui.Prompt{
		Default:   strconv.FormatInt(int64(cfg.End), 10),
		AllowEdit: true,
		Label:     "End key code",
		Validate: func(input string) error {
			mu.Lock()
			defer mu.Unlock()

			i, err := strconv.ParseInt(input, 10, 16)
			if err != nil {
				return err
			}
			cfg.End = int(i)
			if cfg.End < cfg.Start {
				return fmt.Errorf("end is smaller than start")
			}
			return nil
		},
	}

	_, err = prompt.Run()
	if err != nil {
		return cfg, fmt.Errorf("failed to prompt end key code: %w", err)
	}

	var result string
	toBool := func(input string) bool {
		switch strings.ToLower(input) {
		case "y":
			return true
		case "n", "":
			return false
		default:
			panic(fmt.Sprintf("invalid selection: %s", input))
		}
	}

	guard := func(string) error {
		mu.Lock()
		defer mu.Unlock()
		return nil
	}

	prompt = promptui.Prompt{
		Default:   "N",
		IsConfirm: true,
		Label:     "Shift Modifier",
	}
	result, err = prompt.Run()
	if err != nil && !errors.Is(err, promptui.ErrAbort) {
		return cfg, fmt.Errorf("failed to prompt shift modifier: %w", err)
	}
	cfg.Template.ModifierSHIFT = toBool(result)

	prompt = promptui.Prompt{
		Default:   "N",
		IsConfirm: true,
		Label:     "Alt Modifier",
		Validate:  guard,
	}
	result, err = prompt.Run()
	if err != nil && !errors.Is(err, promptui.ErrAbort) {
		return cfg, fmt.Errorf("failed to prompt alt modifier: %w", err)
	}
	cfg.Template.ModifierALT = toBool(result)

	prompt = promptui.Prompt{
		Default:   "N",
		IsConfirm: true,
		Label:     "Ctrl Modifier",
		Validate:  guard,
	}
	result, err = prompt.Run()
	if err != nil && !errors.Is(err, promptui.ErrAbort) {
		return cfg, fmt.Errorf("failed to prompt ctrl modifier: %w", err)
	}
	cfg.Template.ModifierCTRL = toBool(result)

	prompt = promptui.Prompt{
		Default:   "N",
		IsConfirm: true,
		Label:     "Super Modifier",
		Validate:  guard,
	}
	result, err = prompt.Run()
	if err != nil && !errors.Is(err, promptui.ErrAbort) {
		return cfg, fmt.Errorf("failed to prompt win/cmd/super modifier: %w", err)
	}
	cfg.Template.ModifierSuper = toBool(result)

	return cfg, nil
}

func promptRune() (r rune, err error) {
	prompt := promptui.Prompt{
		Label: "Which character were you looking for?",
		Validate: func(s string) error {
			runes := []rune(s)
			if len(runes) != 1 {
				return errors.New("please provide exactly one unicode character")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("failed to prompt single character: %w", err)
	}

	return []rune(result)[0], nil
}

func promptContinue(question string) bool {
	prompt := promptui.Prompt{
		Default:   "Y",
		IsConfirm: true,
		Label:     question,
	}
	_, err := prompt.Run()
	return err == nil
}
