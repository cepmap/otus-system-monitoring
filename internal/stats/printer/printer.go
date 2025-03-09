package printer

import (
	"fmt"

	"github.com/cepmap/otus-system-monitoring/internal/stats/cpu"
	"github.com/cepmap/otus-system-monitoring/internal/stats/diskStat"
	"github.com/cepmap/otus-system-monitoring/internal/stats/disksLoad"
	"github.com/cepmap/otus-system-monitoring/internal/stats/loadAvg"
)

func PrintStats() {
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
}
