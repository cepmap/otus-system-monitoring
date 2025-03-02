package cpu

import (
	"github.com/cepmap/otus-system-monitoring/internal/models"
)

func GetStats() (*models.CPUStat, error) {
	cpuInfo, err := GetCpuStat()
	return cpuInfo, err
}
