package diskLoad

import (
	"github.com/cepmap/otus-system-monitoring/internal/models"
)

func GetStats() (*models.DiskLoad, error) {
	diskLoad, err := GetDiskLoad()
	return diskLoad, err
}
