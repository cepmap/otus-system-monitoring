package integration

import (
	"testing"
	"time"

	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"github.com/cepmap/otus-system-monitoring/internal/collector"
	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/cepmap/otus-system-monitoring/internal/converter"
	"github.com/cepmap/otus-system-monitoring/internal/metrics"
	"github.com/cepmap/otus-system-monitoring/internal/models"
	"github.com/stretchr/testify/require"
)

func TestMetricsIntegration(t *testing.T) {

	config.DaemonConfig = &config.Config{}
	config.DaemonConfig.Stats.LoadAverage = true
	config.DaemonConfig.Stats.Cpu = true
	config.DaemonConfig.Stats.DiskInfo = true
	config.DaemonConfig.Stats.DiskLoad = true

	t.Run("full metrics pipeline", func(t *testing.T) {

		storage := metrics.New()
		require.NotNil(t, storage)

		now := time.Now()

		loadAvg := &models.LoadAverage{Load1Min: 1.0, Load5Min: 2.0, Load15Min: 3.0}
		storage.StoreLoadAverage(loadAvg, now)

		cpuStats := &models.CPUStat{User: 10.0, System: 20.0, Idle: 70.0}
		storage.StoreCPUStats(cpuStats, now)

		disksLoad := &models.DisksLoad{
			DisksLoad: []models.DiskLoad{
				{FSName: "sda1", Tps: 10.0, Kps: 100.0},
			},
		}
		storage.StoreDisksLoad(disksLoad, now)

		diskStats := &models.DiskStats{
			DiskStats: []models.DiskStat{
				{
					FileSystem: "/dev/sda1",
					Usage: models.DiskUsage{
						Used:  500000,
						Usage: "50%",
					},
					Inodes: models.InodeUsage{
						Used:  1000,
						Usage: "10%",
					},
				},
			},
		}
		storage.StoreDiskUsage(diskStats, now)

		statTypes := []pb.StatType{
			pb.StatType_LOAD_AVERAGE,
			pb.StatType_CPU_STATS,
			pb.StatType_DISKS_LOAD,
			pb.StatType_DISK_USAGE,
		}
		avgPeriod := 5 * time.Second
		col := collector.New(storage, statTypes, avgPeriod)
		require.NotNil(t, col)

		protoLoadAvg := converter.LoadAverageToProto(loadAvg)
		require.NotNil(t, protoLoadAvg)
		require.Equal(t, loadAvg.Load1Min, protoLoadAvg.Load1Min)
		require.Equal(t, loadAvg.Load5Min, protoLoadAvg.Load5Min)
		require.Equal(t, loadAvg.Load15Min, protoLoadAvg.Load15Min)

		protoCPUStats := converter.CPUStatToProto(cpuStats)
		require.NotNil(t, protoCPUStats)
		require.Equal(t, cpuStats.User, protoCPUStats.User)
		require.Equal(t, cpuStats.System, protoCPUStats.System)
		require.Equal(t, cpuStats.Idle, protoCPUStats.Idle)

		protoDisksLoad := converter.DisksLoadToProto(disksLoad)
		require.NotNil(t, protoDisksLoad)
		require.Len(t, protoDisksLoad.DisksLoad, 1)
		require.Equal(t, disksLoad.DisksLoad[0].FSName, protoDisksLoad.DisksLoad[0].FsName)
		require.Equal(t, disksLoad.DisksLoad[0].Tps, protoDisksLoad.DisksLoad[0].Tps)
		require.Equal(t, disksLoad.DisksLoad[0].Kps, protoDisksLoad.DisksLoad[0].Kps)

		protoDiskStats := converter.DiskStatsToProto(diskStats)
		require.NotNil(t, protoDiskStats)
		require.Len(t, protoDiskStats.DiskStats, 1)
		require.Equal(t, diskStats.DiskStats[0].FileSystem, protoDiskStats.DiskStats[0].Filesystem)
		require.Equal(t, diskStats.DiskStats[0].Usage.Used, protoDiskStats.DiskStats[0].Usage.Used)
		require.Equal(t, diskStats.DiskStats[0].Usage.Usage, protoDiskStats.DiskStats[0].Usage.Usage)
		require.Equal(t, diskStats.DiskStats[0].Inodes.Used, protoDiskStats.DiskStats[0].Inodes.Used)
		require.Equal(t, diskStats.DiskStats[0].Inodes.Usage, protoDiskStats.DiskStats[0].Inodes.Usage)

		col.CollectMetrics(now)
		response := col.PrepareResponse()
		require.NotNil(t, response)
		require.NotZero(t, response.GetTimestamp())

		require.NotNil(t, response.GetLoadAverage())
		require.NotNil(t, response.GetCpuStats())
		require.NotNil(t, response.GetDisksLoad())
		require.NotNil(t, response.GetDiskStats())
	})

	t.Run("metrics pipeline with partial data", func(t *testing.T) {

		storage := metrics.New()
		require.NotNil(t, storage)

		now := time.Now()

		cpuStats := &models.CPUStat{User: 10.0, System: 20.0, Idle: 70.0}
		storage.StoreCPUStats(cpuStats, now)

		statTypes := []pb.StatType{
			pb.StatType_CPU_STATS,
		}
		avgPeriod := 5 * time.Second
		col := collector.New(storage, statTypes, avgPeriod)
		require.NotNil(t, col)

		protoCPUStats := converter.CPUStatToProto(cpuStats)
		require.NotNil(t, protoCPUStats)
		require.Equal(t, cpuStats.User, protoCPUStats.User)
		require.Equal(t, cpuStats.System, protoCPUStats.System)
		require.Equal(t, cpuStats.Idle, protoCPUStats.Idle)

		col.CollectMetrics(now)
		response := col.PrepareResponse()
		require.NotNil(t, response)
		require.NotZero(t, response.GetTimestamp())

		require.NotNil(t, response.GetCpuStats())
		require.Nil(t, response.GetLoadAverage())
		require.Nil(t, response.GetDisksLoad())
		require.Nil(t, response.GetDiskStats())
	})

	t.Run("metrics pipeline with averaging", func(t *testing.T) {

		storage := metrics.New()
		require.NotNil(t, storage)

		now := time.Now()
		avgPeriod := 3 * time.Second

		stats1 := &models.CPUStat{User: 10.0, System: 20.0, Idle: 70.0}
		stats2 := &models.CPUStat{User: 20.0, System: 30.0, Idle: 50.0}
		stats3 := &models.CPUStat{User: 30.0, System: 40.0, Idle: 30.0}

		storage.StoreCPUStats(stats1, now.Add(-avgPeriod+time.Second))
		storage.StoreCPUStats(stats2, now.Add(-avgPeriod+2*time.Second))
		storage.StoreCPUStats(stats3, now)

		statTypes := []pb.StatType{
			pb.StatType_CPU_STATS,
		}
		col := collector.New(storage, statTypes, avgPeriod)
		require.NotNil(t, col)

		col.CollectMetrics(now)
		response := col.PrepareResponse()
		require.NotNil(t, response)
		require.NotZero(t, response.GetTimestamp())

		avgCPU := response.GetCpuStats()
		require.NotNil(t, avgCPU)
		require.InDelta(t, 20.0, avgCPU.User, 20)
		require.InDelta(t, 30.0, avgCPU.System, 30)
		require.InDelta(t, 50.0, avgCPU.Idle, 50)
	})
}
