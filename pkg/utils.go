package pkg

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func InitControlRequest() (runtime.MemStats, time.Time) {
	var memStatus runtime.MemStats
	start_time := time.Now()

	return memStatus, start_time
}

func FinishControlCheck(memStatus runtime.MemStats, start_time time.Time) (int64, string) {
	runtime.ReadMemStats(&memStatus)
	duration := time.Since(start_time)

	milliseconds := duration.Milliseconds()
	info := fmt.Sprintf("Elapsed time = %vms. Total memory(KB) consumed = %v", milliseconds, memStatus.Sys/1024)

	return milliseconds, info
}

func RoundFloat(num float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(num*p) / p
}
