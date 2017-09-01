package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var buffer = bytes.NewBufferString("")

type Graph struct {
	min    float64
	max    float64
	values []float64
}

func NewGraph(width uint64) *Graph {
	graph := new(Graph)
	graph.min = 0.0
	graph.max = 1.0
	graph.values = make([]float64, width)
	return graph
}

func (graph *Graph) AddValue(val float64) {
	graph.values = append(graph.values[1:], []float64{val}...)
}

func (graph *Graph) Call() {
	diff := graph.max - graph.min
	for _, element := range graph.values {
		ratio := (element - graph.min) / diff
		if ratio <= 0.125 {
			buffer.WriteRune('▁')
			continue
		}
		if ratio <= 0.25 {
			buffer.WriteRune('▂')
			continue
		}
		if ratio <= 0.375 {
			buffer.WriteRune('▃')
			continue
		}
		if ratio <= 0.5 {
			buffer.WriteRune('▄')
			continue
		}
		if ratio <= 0.625 {
			buffer.WriteRune('▅')
			continue
		}
		if ratio <= 0.75 {
			buffer.WriteRune('▆')
			continue
		}
		if ratio <= 0.875 {
			buffer.WriteRune('▇')
			continue
		}
		buffer.WriteRune('█')
	}
}

type Cpu struct {
	idle    uint64
	cpuTime uint64
	graph   *Graph
}

func NewCpu() *Cpu {
	cpu := new(Cpu)
	cpu.graph = NewGraph(10)
	cpu.idle, cpu.cpuTime = cpu.ReadData()
	return cpu
}

func (cpu *Cpu) ReadData() (uint64, uint64) {
	data, _ := ioutil.ReadFile("/proc/stat")
	fields := strings.Fields(strings.Split(string(data), "\n")[0])

	user, _ := strconv.ParseUint(fields[0], 10, 64)
	nice, _ := strconv.ParseUint(fields[1], 10, 64)
	system, _ := strconv.ParseUint(fields[2], 10, 64)
	idle, _ := strconv.ParseUint(fields[3], 10, 64)
	iowait, _ := strconv.ParseUint(fields[4], 10, 64)
	irq, _ := strconv.ParseUint(fields[5], 10, 64)
	softirq, _ := strconv.ParseUint(fields[6], 10, 64)
	steal, _ := strconv.ParseUint(fields[7], 10, 64)
	guest, _ := strconv.ParseUint(fields[8], 10, 64)
	guestNice, _ := strconv.ParseUint(fields[9], 10, 64)
	cpuTime := user + nice + system + idle + iowait + irq + softirq + steal + guest + guestNice

	return idle, cpuTime
}

func (cpu *Cpu) Call() {

	idle, cpuTime := cpu.ReadData()

	fmt.Println(cpuTime)

	idleDiff := idle - cpu.idle
	cpuTimeDiff := cpuTime - cpu.cpuTime

	if cpuTimeDiff == 0 {
		buffer.WriteString("{\"full_text\": \"CPU:\"}")
		return
	}

	usage := float64(idleDiff) / float64(cpuTimeDiff)

	buffer.WriteString("{\"full_text\": \"CPU: ")
	buffer.WriteString(strconv.FormatFloat(usage*100, 'f', 1, 64))
	buffer.WriteString("%\"}")

	cpu.idle = idle
	cpu.cpuTime = cpuTime
}

type Memory struct {
	graph *Graph
}

func NewMemory() *Memory {
	return &Memory{NewGraph(10)}
}

func (mem *Memory) Call() {
	const kB2GB = 1024. * 1024.

	data, _ := ioutil.ReadFile("/proc/meminfo")
	lines := strings.Split(string(data), "\n")

	var total uint64
	var available uint64
	var swap_used uint64
	var swap_total uint64

	total_found := false
	available_found := false
	swap_used_found := false
	swap_total_found := false

	for _, line := range lines {
		if !total_found && strings.HasPrefix(line, "MemTotal") {
			total, _ = strconv.ParseUint(strings.Fields(line)[1], 10, 64)
			total_found = true
		}

		if !available_found && strings.HasPrefix(line, "MemAvailable") {
			available, _ = strconv.ParseUint(strings.Fields(line)[1], 10, 64)
			available_found = true
		}

		if !swap_used_found && strings.HasPrefix(line, "SwapCached") {
			swap_used, _ = strconv.ParseUint(strings.Fields(line)[1], 10, 64)
			swap_used_found = true
		}

		if !swap_total_found && strings.HasPrefix(line, "SwapTotal") {
			swap_total, _ = strconv.ParseUint(strings.Fields(line)[1], 10, 64)
			swap_total_found = true
		}
	}

	mem_usage := 1. - float64(available)/float64(total)
	mem.graph.AddValue(mem_usage)

	buffer.WriteString("{\"full_text\": \"")
	buffer.WriteString("MEM: ")
	buffer.WriteString(strconv.FormatFloat(float64(total-available)/kB2GB, 'f', 1, 64))
	buffer.WriteRune('/')
	buffer.WriteString(strconv.FormatFloat(float64(total)/kB2GB, 'f', 1, 64))
	buffer.WriteString("GB (")
	buffer.WriteString(strconv.FormatFloat(mem_usage*100., 'f', 2, 64))
	buffer.WriteString("%) ")
	mem.graph.Call()
	if swap_total > 0 {
		buffer.WriteString("    SWP: ")
		buffer.WriteString(strconv.FormatFloat(float64(swap_used)/kB2GB, 'f', 1, 64))
		buffer.WriteRune('/')
		buffer.WriteString(strconv.FormatFloat(float64(swap_total)/kB2GB, 'f', 1, 64))
		buffer.WriteString("GB (")
		buffer.WriteString(strconv.FormatFloat(float64(swap_used)/float64(swap_total)*100., 'f', 2, 64))
		buffer.WriteString("%) \"}")
	} else {
		buffer.WriteString(" \"}")
	}
}

type NetworkInterface struct {
	name          string
	lastTime      time.Time
	lastRxBytes   uint64
	lastTxBytes   uint64
	maxRxSpeed    uint64
	maxTxSpeed    uint64
	maxRxSpeedStr string
	maxTxSpeedStr string
}

func NewNetworkInterface(name string) *NetworkInterface {
	net := new(NetworkInterface)

	net.name = name
	net.lastTime = time.Now()
	net.lastRxBytes = net.GetRxBytes()
	net.lastTxBytes = net.GetTxBytes()
	net.maxRxSpeed = 0
	net.maxTxSpeed = 0
	net.maxRxSpeedStr = "0"
	net.maxTxSpeedStr = "0"

	return net
}

func (net *NetworkInterface) GetRxBytes() uint64 {
	data, _ := ioutil.ReadFile("/sys/class/net/" + net.name + "/statistics/rx_bytes")
	rxBytes, _ := strconv.ParseUint(string(data[:len(data)-1]), 10, 64)
	return rxBytes
}

func (net *NetworkInterface) GetTxBytes() uint64 {
	data, _ := ioutil.ReadFile("/sys/class/net/" + net.name + "/statistics/tx_bytes")
	txBytes, _ := strconv.ParseUint(string(data[:len(data)-1]), 10, 64)
	return txBytes
}

func (net *NetworkInterface) Call() {
	currentTime := time.Now()
	elapsedTime := float64(currentTime.Sub(net.lastTime)) / float64(time.Second)

	rxBytes := net.GetRxBytes()
	txBytes := net.GetTxBytes()

	rxSpeed := uint64(float64(rxBytes-net.lastRxBytes) / (elapsedTime * 1024.))
	txSpeed := uint64(float64(txBytes-net.lastTxBytes) / (elapsedTime * 1024.))

	net.lastRxBytes = rxBytes
	net.lastTxBytes = txBytes
	net.lastTime = currentTime

	rxSpeedStr := strconv.FormatUint(rxSpeed, 10)
	txSpeedStr := strconv.FormatUint(txSpeed, 10)

	if rxSpeed > net.maxRxSpeed {
		net.maxRxSpeed = rxSpeed
		net.maxRxSpeedStr = rxSpeedStr
	}

	if txSpeed > net.maxTxSpeed {
		net.maxTxSpeed = txSpeed
		net.maxTxSpeedStr = txSpeedStr
	}

	buffer.WriteString("{\"full_text\": \"")
	buffer.WriteString(net.name)
	buffer.WriteString(": ")
	buffer.WriteRune('↧')
	buffer.WriteString(rxSpeedStr)
	buffer.WriteString("kB/s  ")
	buffer.WriteRune('↥')
	buffer.WriteString(txSpeedStr)
	buffer.WriteString("kB/s\"}")
}

func AddCurrentTime() {
	buffer.WriteString("{\"full_text\": \"")
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	buffer.WriteString("\"}")
}

func AddSeparator() {
	buffer.WriteString("{\"full_text\": \"|\"}")
}

func AddDrive(name string, path string) {
	const BtoGB = float64(1024 * 1024 * 1024)

	var stat syscall.Statfs_t
	syscall.Statfs(path, &stat)
	total := float64(stat.Blocks*uint64(stat.Bsize)) / BtoGB
	used := float64((stat.Blocks-stat.Bfree)*uint64(stat.Bsize)) / BtoGB
	total_str := strconv.FormatFloat(total, 'f', 1, 64)
	used_str := strconv.FormatFloat(used, 'f', 1, 64)
	ratio := used / total * 100.

	buffer.WriteString("{\"full_text\": \"")
	buffer.WriteString(name)
	buffer.WriteString(": ")
	buffer.WriteString(used_str)
	buffer.WriteRune('/')
	buffer.WriteString(total_str)
	buffer.WriteString("GB (")
	buffer.WriteString(strconv.FormatFloat(ratio, 'f', 2, 64))
	buffer.WriteString("%)")
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

	// cpu := NewCpu()
	mem := NewMemory()
	net_lan := NewNetworkInterface("enp8s0")

	for true {
		BeginStatus()
		// cpu.Call()
		// Comma()
		mem.Call()
		Comma()
		AddDrive("SSD", "/")
		Comma()
		net_lan.Call()
		Comma()
		AddCurrentTime()
		EndStatus()
	}
}
