package collector

import (
	"testing"
	"time"

	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/cepmap/otus-system-monitoring/internal/metrics"
	"github.com/stretchr/testify/require"
)

func TestCollector(t *testing.T) {
	config.DaemonConfig = &config.Config{}
	config.DaemonConfig.Stats.LoadAverage = true
	config.DaemonConfig.Stats.Cpu = true
	config.DaemonConfig.Stats.DiskInfo = true
	config.DaemonConfig.Stats.DiskLoad = true

	t.Run("create collector", func(t *testing.T) {
		metricsStorage := metrics.New()
		statTypes := []pb.StatType{
			pb.StatType_LOAD_AVERAGE,
			pb.StatType_CPU_STATS,
			pb.StatType_DISKS_LOAD,
			pb.StatType_DISK_USAGE,
		}
		avgPeriod := 5 * time.Second

		collector := New(metricsStorage, statTypes, avgPeriod)
		require.NotNil(t, collector)
		require.Equal(t, metricsStorage, collector.metrics)
		require.Equal(t, statTypes, collector.statTypes)
		require.Equal(t, avgPeriod, collector.avgPeriod)
	})

	t.Run("collect metrics", func(t *testing.T) {
		metricsStorage := metrics.New()
		statTypes := []pb.StatType{
			pb.StatType_LOAD_AVERAGE,
			pb.StatType_CPU_STATS,
			pb.StatType_DISKS_LOAD,
			pb.StatType_DISK_USAGE,
		}
		avgPeriod := 5 * time.Second

		collector := New(metricsStorage, statTypes, avgPeriod)
		require.NotNil(t, collector)

		collector.CollectInitialData()

		timestamp := time.Now()
		collector.CollectMetrics(timestamp)

		response := collector.PrepareResponse()
		require.NotNil(t, response)
		require.NotZero(t, response.GetTimestamp())

		if config.DaemonConfig.Stats.LoadAverage {
			require.NotNil(t, response.GetLoadAverage())
		}
		if config.DaemonConfig.Stats.Cpu {
			require.NotNil(t, response.GetCpuStats())
		}
		if config.DaemonConfig.Stats.DiskInfo {
			require.NotNil(t, response.GetDiskStats())
		}
		if config.DaemonConfig.Stats.DiskLoad {
			require.NotNil(t, response.GetDisksLoad())
		}
	})

	t.Run("collect specific metrics", func(t *testing.T) {
		metricsStorage := metrics.New()

		statTypes := []pb.StatType{
			pb.StatType_CPU_STATS,
		}
		avgPeriod := 5 * time.Second

		collector := New(metricsStorage, statTypes, avgPeriod)
		require.NotNil(t, collector)

		collector.CollectInitialData()

		timestamp := time.Now()
		collector.CollectMetrics(timestamp)

		response := collector.PrepareResponse()
		require.NotNil(t, response)
		require.NotZero(t, response.GetTimestamp())

		require.NotNil(t, response.GetCpuStats())
		require.Nil(t, response.GetLoadAverage())
		require.Nil(t, response.GetDiskStats())
		require.Nil(t, response.GetDisksLoad())
	})
}
