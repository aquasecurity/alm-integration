package msgservice

import "log"

func prnInputLogs(msg string, v ...interface{}) {
	maxLen := 20
	for idx, e := range v {
		b, ok := e.([]byte)
		if ok {
			if l := len(b); l > maxLen {
				v[idx] = string(b[:maxLen])
			}
		}
	}
	log.Printf(msg, v...)
}
