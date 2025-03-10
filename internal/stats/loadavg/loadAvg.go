package loadavg

import "github.com/cepmap/otus-system-monitoring/internal/models"

func GetStats() (*models.LoadAverage, error) {
	loadAvg, err := GetStatsOs()

	return loadAvg, err
}
