package main

import gogobar "github.com/f-koehler/gogobar/gogobar"

func main() {
	gogobar.Init(500)

	cpu := gogobar.NewCpu()
	drive_ssd := gogobar.NewDrive("SSD", "/")
	mem := gogobar.NewMemory()
	net_lan := gogobar.NewNetworkInterface("enp8s0")

	for true {
		gogobar.BeginStatus()
		cpu.Call()
		gogobar.Comma()
		mem.Call()
		gogobar.Comma()
		drive_ssd.Call()
		gogobar.Comma()
		net_lan.Call()
		gogobar.Comma()
		gogobar.CurrentTime()
		gogobar.EndStatus()
	}
}
