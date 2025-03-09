package main

import (
	"fmt"

	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/cepmap/otus-system-monitoring/internal/logger"
	"github.com/cepmap/otus-system-monitoring/internal/stats/cpu"
	"github.com/cepmap/otus-system-monitoring/internal/stats/diskStat"
	"github.com/cepmap/otus-system-monitoring/internal/stats/disksLoad"
	"github.com/cepmap/otus-system-monitoring/internal/stats/loadAvg"
)

func main() {
	logger.Info(fmt.Sprintf("Current config: %+v", config.DaemonConfig))
	res, err := disksLoad.GetStats()
	if err != nil {
		return
	}

	fmt.Println(res)

	res1, err := cpu.GetCpuStat()
	if err != nil {
		return
	}
	fmt.Println(res1)

	res2, err := loadAvg.GetStats()
	if err != nil {
		return
	}
	fmt.Println(res2)

	res3, err := diskStat.GetStats()
	if err != nil {
		return
	}
	fmt.Println(res3)

	//ticker := time.NewTicker(1 * time.Second)
	//defer ticker.Stop() // Ensure the ticker is stopped when done
	//
	//// Run a loop that executes a function every second
	//for range ticker.C {
	//	// Call your function here
	//	res, err := diskStat.GetStats()
	//	if err != nil {
	//		return
	//	}
	//	fmt.Println(res)
	//}
}
