package printer

import (
	"fmt"

	"github.com/cepmap/otus-system-monitoring/internal/stats/cpu"
	"github.com/cepmap/otus-system-monitoring/internal/stats/disksload"
	"github.com/cepmap/otus-system-monitoring/internal/stats/diskstat"
	"github.com/cepmap/otus-system-monitoring/internal/stats/loadavg"
)

func PrintStats() {
	res, err := disksload.GetStats()
	if err != nil {
		return
	}
	fmt.Println(res)

	res1, err := cpu.GetCpuStat()
	if err != nil {
		return
	}
	fmt.Println(res1)

	res2, err := loadavg.GetStats()
	if err != nil {
		return
	}
	fmt.Println(res2)

	res3, err := diskstat.GetStats()
	if err != nil {
		return
	}
	fmt.Println(res3)
}
