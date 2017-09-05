package gogobar

import (
	"io/ioutil"
	"strconv"
	"strings"
)

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
	var swapUsed uint64
	var swapTotal uint64

	totalFound := false
	availableFound := false
	swapUsedFound := false
	swapTotalFound := false

	for _, line := range lines {
		if !totalFound && strings.HasPrefix(line, "MemTotal") {
			total, _ = strconv.ParseUint(strings.Fields(line)[1], 10, 64)
			totalFound = true
		}

		if !availableFound && strings.HasPrefix(line, "MemAvailable") {
			available, _ = strconv.ParseUint(strings.Fields(line)[1], 10, 64)
			availableFound = true
		}

		if !swapUsedFound && strings.HasPrefix(line, "SwapCached") {
			swapUsed, _ = strconv.ParseUint(strings.Fields(line)[1], 10, 64)
			swapUsedFound = true
		}

		if !swapTotalFound && strings.HasPrefix(line, "SwapTotal") {
			swapTotal, _ = strconv.ParseUint(strings.Fields(line)[1], 10, 64)
			swapTotalFound = true
		}
	}

	memUsage := 1. - float64(available)/float64(total)
	mem.graph.AddValue(memUsage)

	usedStr := strconv.FormatFloat(float64(total-available)/kB2GB, 'f', 1, 64)
	totalStr := strconv.FormatFloat(float64(total)/kB2GB, 'f', 1, 64)
	usageStr := strconv.FormatFloat(memUsage*100., 'f', 2, 64)

	buffer.WriteString("{\"full_text\": \"")
	buffer.WriteString("MEM: ")
	buffer.WriteString(PadLeft(usedStr, ' ', len(totalStr)))
	buffer.WriteRune('/')
	buffer.WriteString(totalStr)
	buffer.WriteString("GB (")
	buffer.WriteString(PadLeft(usageStr, ' ', 6))
	buffer.WriteString("%) ")
	mem.graph.Call()

	if swapTotal > 0 {
		usedStr = strconv.FormatFloat(float64(swapUsed)/kB2GB, 'f', 1, 64)
		totalStr = strconv.FormatFloat(float64(swapTotal)/kB2GB, 'f', 1, 64)
		usageStr = strconv.FormatFloat(float64(swapUsed)/float64(swapTotal)*100., 'f', 2, 64)

		buffer.WriteString("    SWP: ")
		buffer.WriteString(PadLeft(usedStr, ' ', len(totalStr)))
		buffer.WriteRune('/')
		buffer.WriteString(totalStr)
		buffer.WriteString("GB (")
		buffer.WriteString(PadLeft(usageStr, ' ', 6))
		buffer.WriteString("%) \"}")
	} else {
		buffer.WriteString(" \"}")
	}
}
