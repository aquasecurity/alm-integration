package scanservice

import "log"

func prnInputLogs(msg string, input []byte, v ...interface{}) {
	maxLen := 100
	if l := len(input); l < maxLen {
		maxLen = l
	}
	log.Printf(msg, string(input[:maxLen]), v)
}
