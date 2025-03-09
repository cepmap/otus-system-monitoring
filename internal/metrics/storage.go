package metrics

import (
	"time"

	"github.com/cepmap/otus-system-monitoring/internal/models"
	"github.com/cepmap/otus-system-monitoring/internal/storage"
	memorystorage "github.com/cepmap/otus-system-monitoring/internal/storage/memory"
)

type MetricsStorage struct {
	loadAvg   storage.Storage
	cpuStats  storage.Storage
	diskLoad  storage.Storage
	diskUsage storage.Storage
}

func New() *MetricsStorage {
	return &MetricsStorage{
		loadAvg:   memorystorage.New(),
		cpuStats:  memorystorage.New(),
		diskLoad:  memorystorage.New(),
		diskUsage: memorystorage.New(),
	}
}

func (m *MetricsStorage) StoreLoadAverage(stats *models.LoadAverage, timestamp time.Time) {
	m.loadAvg.Push(stats, timestamp)
}

func (m *MetricsStorage) StoreCPUStats(stats *models.CPUStat, timestamp time.Time) {
	m.cpuStats.Push(stats, timestamp)
}

func (m *MetricsStorage) StoreDisksLoad(stats *models.DisksLoad, timestamp time.Time) {
	m.diskLoad.Push(stats, timestamp)
}

func (m *MetricsStorage) StoreDiskUsage(stats *models.DiskStats, timestamp time.Time) {
	m.diskUsage.Push(stats, timestamp)
}

func getAverageFromStorage[T any](store storage.Storage, period time.Duration) []T {
	now := time.Now()
	start := now.Add(-period)

	var result []T
	for item := range store.GetElementsAt(start) {
		if stat, ok := item.(T); ok {
			result = append(result, stat)
		}
	}
	return result
}

func (m *MetricsStorage) GetAverageLoadAverage(period time.Duration) *models.LoadAverage {
	stats := getAverageFromStorage[*models.LoadAverage](m.loadAvg, period)
	return averageLoadAverage(stats)
}

func (m *MetricsStorage) GetAverageCPUStats(period time.Duration) *models.CPUStat {
	stats := getAverageFromStorage[*models.CPUStat](m.cpuStats, period)
	return averageCPUStat(stats)
}

func (m *MetricsStorage) GetAverageDisksLoad(period time.Duration) *models.DisksLoad {
	stats := getAverageFromStorage[*models.DisksLoad](m.diskLoad, period)
	return averageDisksLoad(stats)
}

func (m *MetricsStorage) GetLatestDiskUsage() *models.DiskStats {
	stats := getAverageFromStorage[*models.DiskStats](m.diskUsage, time.Second)
	if len(stats) > 0 {
		return stats[len(stats)-1]
	}
	return nil
}
