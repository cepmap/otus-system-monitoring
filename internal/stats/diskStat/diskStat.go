package diskStat

import (
	"github.com/cepmap/otus-system-monitoring/internal/models"
)

func GetStats() (*models.DiskStats, error) {
	diskStat, err := GetDiskStats()
	return diskStat, err
}
