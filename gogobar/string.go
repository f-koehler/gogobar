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

func PadRight(str string, pad rune, length int) string {
	currentLen := len(str)
	if currentLen >= length {
		return str
	}

	diff := length - currentLen
	return str + strings.Repeat(string(pad), diff)
}

func PadBoth(str string, pad rune, length int) string {
	currentLen := len(str)
	if currentLen >= length {
		return str
	}

	diff := length - currentLen
	left := diff / 2
	right := diff - left
	return strings.Repeat(string(pad), left) + str + strings.Repeat(string(pad), right)
}
