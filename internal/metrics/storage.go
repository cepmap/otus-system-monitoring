package metrics

import (
	"sync"
	"time"

	"github.com/cepmap/otus-system-monitoring/internal/models"
	"github.com/cepmap/otus-system-monitoring/internal/storage"
	memorystorage "github.com/cepmap/otus-system-monitoring/internal/storage/memory"
)

type Storage struct {
	mu        sync.RWMutex
	loadAvg   storage.Storage
	cpuStats  storage.Storage
	diskLoad  storage.Storage
	diskUsage storage.Storage
}

func New() *Storage {
	return &Storage{
		loadAvg:   memorystorage.New(),
		cpuStats:  memorystorage.New(),
		diskLoad:  memorystorage.New(),
		diskUsage: memorystorage.New(),
	}
}

func (m *Storage) StoreLoadAverage(stats *models.LoadAverage, timestamp time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loadAvg.Push(stats, timestamp)
}

func (m *Storage) StoreCPUStats(stats *models.CPUStat, timestamp time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cpuStats.Push(stats, timestamp)
}

func (m *Storage) StoreDisksLoad(stats *models.DisksLoad, timestamp time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.diskLoad.Push(stats, timestamp)
}

func (m *Storage) StoreDiskUsage(stats *models.DiskStats, timestamp time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
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

func (m *Storage) GetAverageLoadAverage(period time.Duration) *models.LoadAverage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	stats := getAverageFromStorage[*models.LoadAverage](m.loadAvg, period)
	return averageLoadAverage(stats)
}

func (m *Storage) GetAverageCPUStats(period time.Duration) *models.CPUStat {
	m.mu.RLock()
	defer m.mu.RUnlock()
	stats := getAverageFromStorage[*models.CPUStat](m.cpuStats, period)
	return averageCPUStat(stats)
}

func (m *Storage) GetAverageDisksLoad(period time.Duration) *models.DisksLoad {
	m.mu.RLock()
	defer m.mu.RUnlock()
	stats := getAverageFromStorage[*models.DisksLoad](m.diskLoad, period)
	return averageDisksLoad(stats)
}

func (m *Storage) GetLatestDiskUsage() *models.DiskStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	stats := getAverageFromStorage[*models.DiskStats](m.diskUsage, time.Second)
	if len(stats) > 0 {
		return stats[len(stats)-1]
	}
	return nil
}
