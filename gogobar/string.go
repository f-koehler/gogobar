package gogobar

import (
	"strings"
)

func PadLeft(str string, pad rune, length int) string {
	currentLen := len(str)
	if currentLen >= length {
		return str
	}

	diff := length - currentLen
	return strings.Repeat(string(pad), diff) + str
}
