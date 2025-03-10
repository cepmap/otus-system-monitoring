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

func (m *Storage) StartCleaner(ctx context.Context) {
	go m.cleanerLoop(ctx)
}

func (m *Storage) cleanerLoop(ctx context.Context) {
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

func (m *Storage) cleanOldData() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-defaultRetentionPeriod)

	cleanedCount := 0
	cleanedCount += m.cleanStorageOldData(m.loadAvg, cutoff)
	cleanedCount += m.cleanStorageOldData(m.cpuStats, cutoff)
	cleanedCount += m.cleanStorageOldData(m.diskLoad, cutoff)
	cleanedCount += m.cleanStorageOldData(m.diskUsage, cutoff)

	logger.Info(fmt.Sprintf("Cleaned %d old metrics data before %s", cleanedCount, cutoff.Format(time.RFC3339)))
}

func (m *Storage) cleanStorageOldData(s storage.Storage, cutoff time.Time) int {
	count := 0
	for item := range s.GetElementsAt(time.Time{}) {
		if ts, ok := s.GetTimestamp(item); ok && ts.Before(cutoff) {
			s.Remove(item)
			count++
		}
	}
	return count
}
