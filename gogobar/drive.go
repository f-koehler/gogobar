package gogobar

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
	usage := float64((stat.Blocks-stat.Bfree)*uint64(stat.Bsize)) / BtoGB
	totalStr := strconv.FormatFloat(total, 'f', 1, 64)
	usageStr := strconv.FormatFloat(usage, 'f', 1, 64)
	ratio := usage / total * 100.

	buffer.WriteString("{\"full_text\": \"")
	buffer.WriteString(drive.name)
	buffer.WriteString(": ")
	buffer.WriteString(usageStr)
	buffer.WriteRune('/')
	buffer.WriteString(totalStr)
	buffer.WriteString("GB (")
	buffer.WriteString(strconv.FormatFloat(ratio, 'f', 2, 64))
	buffer.WriteString("%)")

	buffer.WriteString("\", \"color\": \"")
	if usage > 0.75 {
		buffer.WriteString(colorCritical)
	} else if usage > 0.5 {
		buffer.WriteString(colorBad)
	} else if usage > 0.25 {
		buffer.WriteString(colorGood)
	} else {
		buffer.WriteString(colorNeutral)
	}
	buffer.WriteString("\"}")
}
