package sendkeys

import "encoding/json"

type KeyCode struct {
	Code          int  `json:"code"`
	ModifierSuper bool `json:"super"` // WIN/CMD/MOD
	ModifierALT   bool `json:"alt"`   // Alt/Option
	ModifierCTRL  bool `json:"ctrl"`
	ModifierSHIFT bool `json:"shift"`
}

func (k KeyCode) String() string {
	data, _ := json.Marshal(k)
	return string(data)
}

func SimpleKeyCode(code int) KeyCode {
	return KeyCode{
		Code: code,
	}
}

func ShiftKeyCode(code int) KeyCode {
	return KeyCode{
		Code:          code,
		ModifierSHIFT: true,
	}
}

func AltKeyCode(code int) KeyCode {
	return KeyCode{
		Code:        code,
		ModifierALT: true,
	}
}

func AltShiftKeyCode(code int) KeyCode {
	return KeyCode{
		Code:          code,
		ModifierALT:   true,
		ModifierSHIFT: true,
	}
}
