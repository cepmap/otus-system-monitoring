package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/cepmap/otus-system-monitoring/internal/logger"
	"github.com/cepmap/otus-system-monitoring/internal/storage"
)

const (
	defaultCleanupInterval = 5 * time.Minute

	defaultRetentionPeriod = 24 * time.Hour
)

func (m *MetricsStorage) StartCleaner(ctx context.Context) {
	go m.cleanerLoop(ctx)
}

func (m *MetricsStorage) cleanerLoop(ctx context.Context) {
	ticker := time.NewTicker(defaultCleanupInterval)
	defer ticker.Stop()

	logger.Info("Metrics cleaner started")
	defer logger.Info("Metrics cleaner stopped")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.cleanOldData()
		}
	}
}

func (m *MetricsStorage) cleanOldData() {
	now := time.Now()
	cutoff := now.Add(-defaultRetentionPeriod)

	storages := []struct {
		name    string
		storage storage.Storage
	}{
		{"load_average", m.loadAvg},
		{"cpu_stats", m.cpuStats},
		{"disk_load", m.diskLoad},
		{"disk_usage", m.diskUsage},
	}

	for _, s := range storages {
		count := 0
		for item := range s.storage.GetElementsAt(time.Time{}) {
			if ts, ok := s.storage.GetTimestamp(item); ok && ts.Before(cutoff) {
				s.storage.Remove(item)
				count++
			}
		}
		if count > 0 {
			logger.Info(fmt.Sprintf("Cleaned %d old records from %s storage", count, s.name))
		}
	}
}
