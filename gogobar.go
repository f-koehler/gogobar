package main

const refresh_interval = 500

func main() {
	Init()

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
		CurrentTime()
		EndStatus()
	}
}
