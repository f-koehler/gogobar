package main

import (
	"bytes"
	"fmt"
	"time"
)

var buffer = bytes.NewBufferString("")

func AddCurrentTime() {
	buffer.WriteString("{\"full_text\": \"")
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	buffer.WriteString("\"}")
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

func main() {
	fmt.Println("{ \"version\": 1}")
	fmt.Println("[")

	cpu := NewCpu()
	drive_ssd := NewDrive("SSD", "/")
	mem := NewMemory()
	net_lan := NewNetworkInterface("enp8s0")

	for true {
		BeginStatus()
		cpu.Call()
		Comma()
		mem.Call()
		Comma()
		drive_ssd.Call()
		Comma()
		net_lan.Call()
		Comma()
		AddCurrentTime()
		EndStatus()
	}
}
