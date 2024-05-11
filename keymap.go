package sendkeys

import (
	"fmt"
)

func (kb *KBWrap) strToKeys(str string) (keys []keyCode) {
	for _, r := range str {
		code, ok := keyCodes[r]
		if ok {
			keys = append(keys, code)
		} else {
			kb.errors = append(
				kb.errors,
				fmt.Errorf("%w: %v", ErrKeyMappingNotFound, r),
			)
		}
	}
	return
}
