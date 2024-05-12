package sendkeys

func GenerateKeyCodesWithTemplate(tmpl KeyCode, minMax ...int) []KeyCode {
	min := 0
	max := 0x60
	switch len(minMax) {
	case 1:
		max = minMax[0]
	case 2:
		min = minMax[0]
		max = minMax[1]
	}

	if min > max {
		min, max = max, min
	}

	result := make([]KeyCode, 0, max-min)
	for i := min; i < max; i++ {
		kc := tmpl
		kc.Code = i

		result = append(result, kc)
	}
	return result
}
