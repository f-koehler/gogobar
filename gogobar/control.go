package gogobar

import (
	"bytes"
	"fmt"
	"time"
)

var buffer = bytes.NewBufferString("")
var refresh_interval time.Duration

func Init(interval int) {
	refresh_interval = time.Duration(interval)
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
	time.Sleep(refresh_interval * time.Millisecond)
}
