package collector

import (
	"fmt"
	"time"

	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/cepmap/otus-system-monitoring/internal/converter"
	"github.com/cepmap/otus-system-monitoring/internal/logger"
	"github.com/cepmap/otus-system-monitoring/internal/metrics"
	"github.com/cepmap/otus-system-monitoring/internal/stats/cpu"
	"github.com/cepmap/otus-system-monitoring/internal/stats/diskStat"
	"github.com/cepmap/otus-system-monitoring/internal/stats/disksLoad"
	"github.com/cepmap/otus-system-monitoring/internal/stats/loadAvg"
)

type Collector struct {
	metrics   *metrics.MetricsStorage
	statTypes []pb.StatType
	avgPeriod time.Duration
}

func New(metrics *metrics.MetricsStorage, statTypes []pb.StatType, avgPeriod time.Duration) *Collector {
	return &Collector{
		metrics:   metrics,
		statTypes: statTypes,
		avgPeriod: avgPeriod,
	}
}

func (c *Collector) CollectMetrics(timestamp time.Time) {
	for _, statType := range c.statTypes {
		switch statType {
		case pb.StatType_LOAD_AVERAGE:
			if config.DaemonConfig.Stats.LoadAverage {
				if stats, err := loadAvg.GetStats(); err == nil {
					c.metrics.StoreLoadAverage(stats, timestamp)
				}
			}
		case pb.StatType_CPU_STATS:
			if config.DaemonConfig.Stats.Cpu {
				if stats, err := cpu.GetCpuStat(); err == nil {
					c.metrics.StoreCPUStats(stats, timestamp)
				}
			}
		case pb.StatType_DISKS_LOAD:
			if config.DaemonConfig.Stats.DiskLoad {
				if stats, err := disksLoad.GetStats(); err == nil {
					c.metrics.StoreDisksLoad(stats, timestamp)
				}
			}
		case pb.StatType_DISK_USAGE:
			if config.DaemonConfig.Stats.DiskInfo {
				if stats, err := diskStat.GetStats(); err == nil {
					c.metrics.StoreDiskUsage(stats, timestamp)
				}
			}
		}
	}
}

func (c *Collector) CollectInitialData() {
	logger.Info(fmt.Sprintf("Starting initial data collection for %v", c.avgPeriod))
	startTime := time.Now()

	for time.Since(startTime) < c.avgPeriod {
		c.CollectMetrics(time.Now())
		time.Sleep(1 * time.Second)
	}

	logger.Info("Initial data collection completed")
}

func (c *Collector) PrepareResponse() *pb.StatsResponse {
	response := &pb.StatsResponse{
		Timestamp: time.Now().Unix(),
	}

	for _, statType := range c.statTypes {
		switch statType {
		case pb.StatType_LOAD_AVERAGE:
			if config.DaemonConfig.Stats.LoadAverage {
				if avgStats := c.metrics.GetAverageLoadAverage(c.avgPeriod); avgStats != nil {
					response.LoadAverage = converter.LoadAverageToProto(avgStats)
				}
			}
		case pb.StatType_CPU_STATS:
			if config.DaemonConfig.Stats.Cpu {
				if avgStats := c.metrics.GetAverageCPUStats(c.avgPeriod); avgStats != nil {
					response.CpuStats = converter.CPUStatToProto(avgStats)
				}
			}
		case pb.StatType_DISKS_LOAD:
			if config.DaemonConfig.Stats.DiskLoad {
				if avgStats := c.metrics.GetAverageDisksLoad(c.avgPeriod); avgStats != nil {
					response.DisksLoad = converter.DisksLoadToProto(avgStats)
				}
			}
		case pb.StatType_DISK_USAGE:
			if config.DaemonConfig.Stats.DiskInfo {
				if stats := c.metrics.GetLatestDiskUsage(); stats != nil {
					response.DiskStats = converter.DiskStatsToProto(stats)
				}
			}
		}
	}

	return response
}
