package sendkeys

type keyCode struct {
	Code          int
	ModifierSuper bool // WIN/CMD/MOD
	ModifierALT   bool // Alt/Option
	ModifierCTRL  bool
	ModifierSHIFT bool
}

func simpleKeyCode(code int) keyCode {
	return keyCode{
		Code: code,
	}
}

func shiftKeyCode(code int) keyCode {
	return keyCode{
		Code:          code,
		ModifierSHIFT: true,
	}
}

func altKeyCode(code int) keyCode {
	return keyCode{
		Code:        code,
		ModifierALT: true,
	}
}

func altShiftKeyCode(code int) keyCode {
	return keyCode{
		Code:          code,
		ModifierALT:   true,
		ModifierSHIFT: true,
	}
}
