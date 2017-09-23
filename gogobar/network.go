package gogobar

import (
	"io/ioutil"
	"math"
	"strconv"
	"time"
)

type NetworkInterface struct {
	name          string
	lastTime      time.Time
	lastRxBytes   uint64
	lastTxBytes   uint64
	graphRx       *Graph
	graphTx       *Graph
	maxRxSpeed    uint64
	maxTxSpeed    uint64
	maxRxSpeedLen int
	maxTxSpeedLen int
}

func NewNetworkInterface(name string) *NetworkInterface {
	net := new(NetworkInterface)

	net.name = name
	net.lastTime = time.Now()
	net.lastRxBytes = net.GetRxBytes()
	net.lastTxBytes = net.GetTxBytes()
	net.graphRx = NewGraph(10)
	net.graphTx = NewGraph(10)
	net.maxRxSpeed = 0
	net.maxTxSpeed = 0
	net.maxRxSpeedLen = 0
	net.maxTxSpeedLen = 0

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

	rxSpeedStr := strconv.FormatUint(rxSpeed, 10)
	txSpeedStr := strconv.FormatUint(txSpeed, 10)

	net.lastRxBytes = rxBytes
	net.lastTxBytes = txBytes
	net.lastTime = currentTime

	if rxSpeed > net.maxRxSpeed {
		net.maxRxSpeed = rxSpeed
		net.maxRxSpeedLen = len(rxSpeedStr)
		net.graphRx.max = float64(rxSpeed)
	}

	if txSpeed > net.maxTxSpeed {
		net.maxTxSpeed = txSpeed
		net.maxTxSpeedLen = len(txSpeedStr)
		net.graphTx.max = float64(txSpeed)
	}

	usage := math.Max(float64(rxSpeed)/float64(net.maxRxSpeed), float64(txSpeed)/float64(net.maxTxSpeed))

	net.graphRx.AddValue(float64(rxSpeed))
	net.graphTx.AddValue(float64(txSpeed))

	buffer.WriteString("{\"full_text\": \"")
	buffer.WriteString(net.name)
	buffer.WriteString(": ")
	buffer.WriteString(PadLeft("↧"+rxSpeedStr+"kB/s", ' ', net.maxRxSpeedLen+5))
	buffer.WriteRune(' ')
	net.graphRx.Call()
	buffer.WriteString("   ")
	buffer.WriteString(PadLeft("↥"+txSpeedStr+"kB/s", ' ', net.maxTxSpeedLen+5))
	buffer.WriteRune(' ')
	net.graphTx.Call()

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
