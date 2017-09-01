package main

import (
	"strconv"
	"syscall"
)

type Drive struct {
	name string
	path string
}

func NewDrive(name string, path string) *Drive {
	return &Drive{name: name, path: path}
}

func (drive *Drive) Call() {
	const BtoGB = float64(1024 * 1024 * 1024)

	var stat syscall.Statfs_t
	syscall.Statfs(drive.path, &stat)
	total := float64(stat.Blocks*uint64(stat.Bsize)) / BtoGB
	used := float64((stat.Blocks-stat.Bfree)*uint64(stat.Bsize)) / BtoGB
	total_str := strconv.FormatFloat(total, 'f', 1, 64)
	used_str := strconv.FormatFloat(used, 'f', 1, 64)
	ratio := used / total * 100.

	buffer.WriteString("{\"full_text\": \"")
	buffer.WriteString(drive.name)
	buffer.WriteString(": ")
	buffer.WriteString(used_str)
	buffer.WriteRune('/')
	buffer.WriteString(total_str)
	buffer.WriteString("GB (")
	buffer.WriteString(strconv.FormatFloat(ratio, 'f', 2, 64))
	buffer.WriteString("%)")
	buffer.WriteString("\"}")
}
