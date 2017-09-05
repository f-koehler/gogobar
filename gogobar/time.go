package gogobar

import "time"

func CurrentTime() {
	buffer.WriteString("{\"full_text\": \"")
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	buffer.WriteString("\"}")
}
