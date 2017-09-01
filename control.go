package main

import (
	"bytes"
	"fmt"
	"time"
)

var buffer = bytes.NewBufferString("")

func Init() {
	fmt.Println("{ \"version\": 1}")
	fmt.Println("[")
}

func BeginStatus() {
	buffer.Reset()
	buffer.WriteString("[")
}

func Comma() {
	buffer.WriteRune(',')
}

func EndStatus() {
	buffer.WriteString("],")
	fmt.Println(buffer.String())
	buffer.Reset()
	time.Sleep(500 * time.Millisecond)
}
