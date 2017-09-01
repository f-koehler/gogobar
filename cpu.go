package main

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type Cpu struct {
	idle  uint64
	total uint64
	graph *Graph
}

func NewCpu() *Cpu {
	cpu := new(Cpu)
	cpu.graph = NewGraph(10)
	cpu.idle, cpu.total = cpu.ReadData()
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
	total := user + nice + system + idle + iowait + irq + softirq + steal + guest + guestNice

	return idle, total
}

func (cpu *Cpu) Call() {

	idle, total := cpu.ReadData()

	idleDiff := idle - cpu.idle
	totalDiff := total - cpu.total

	if totalDiff == 0 {
		buffer.WriteString("{\"full_text\": \"CPU:\"}")
		return
	}

	usage := float64(idleDiff) / float64(totalDiff)
	cpu.graph.AddValue(usage)

	buffer.WriteString("{\"full_text\": \"CPU: ")
	buffer.WriteString(strconv.FormatFloat(usage*100, 'f', 1, 64))
	buffer.WriteString("% ")
	cpu.graph.Call()
	buffer.WriteString("\"}")

	cpu.idle = idle
	cpu.total = total
}
