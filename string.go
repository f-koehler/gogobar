package main

import (
	"strings"
)

func FillAndWrite(str string, length int, left bool) {
	lengthDiff := length - len(str)
	if lengthDiff <= 0 {
		buffer.WriteString(str)
		return
	}
	if left {
		buffer.WriteString(strings.Repeat(" ", lengthDiff))
		buffer.WriteString(str)
		return
	}
	buffer.WriteString(str)
	buffer.WriteString(strings.Repeat(" ", lengthDiff))
}
