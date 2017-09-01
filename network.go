package main

import (
	"io/ioutil"
	"strconv"
	"time"
)

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
