package main

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
