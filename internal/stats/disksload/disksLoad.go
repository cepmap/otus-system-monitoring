package disksload

import (
	"github.com/cepmap/otus-system-monitoring/internal/models"
)

func GetStats() (*models.DisksLoad, error) {
	diskLoad, err := GetDisksLoad()
	return diskLoad, err
}
