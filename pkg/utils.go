package pkg

import (
	"runtime"
	"time"
)

func InitControlRequest() (runtime.MemStats, time.Time) {
	var memStatus runtime.MemStats
	start_time := time.Now()

	return memStatus, start_time
}

func FinishControlCheck(memStatus runtime.MemStats, start_time time.Time) int64 {
	runtime.ReadMemStats(&memStatus)
	duration := time.Since(start_time)

	milliseconds := duration.Milliseconds()

	return milliseconds
}
