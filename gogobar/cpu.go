package gogobar

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type Cpu struct {
	idle      uint64
	total     uint64
	lastIdle  uint64
	lastTotal uint64
	graph     *Graph
}

func NewCpu() *Cpu {
	cpu := new(Cpu)
	cpu.graph = NewGraph(10)
	cpu.ReadData()
	return cpu
}

func (cpu *Cpu) ReadData() {
	data, _ := ioutil.ReadFile("/proc/stat")
	fields := strings.Fields(strings.Split(string(data), "\n")[0])

	cpu.lastIdle = cpu.idle
	cpu.lastTotal = cpu.total

	cpu.idle, _ = strconv.ParseUint(fields[4], 10, 64)
	cpu.total = 0
	val := uint64(0)
	for i := 1; i < len(fields); i++ {
		val, _ = strconv.ParseUint(fields[i], 10, 64)
		cpu.total += val
	}
}

func (cpu *Cpu) Call() {

	cpu.ReadData()
	idleDiff := cpu.idle - cpu.lastIdle
	totalDiff := cpu.total - cpu.lastTotal

	if totalDiff == 0 {
		buffer.WriteString("{\"full_text\": \"CPU:\"}")
		return
	}

	usage := 1. - float64(idleDiff)/float64(totalDiff)
	cpu.graph.AddValue(usage)

	buffer.WriteString("{\"full_text\": \"CPU: ")
	buffer.WriteString(PadLeft(strconv.FormatFloat(usage*100, 'f', 1, 64), ' ', 5))
	buffer.WriteString("% ")
	cpu.graph.Call()
	buffer.WriteString("\"}")
}
