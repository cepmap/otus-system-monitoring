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
	"github.com/cepmap/otus-system-monitoring/internal/stats/disksload"
	"github.com/cepmap/otus-system-monitoring/internal/stats/diskstat"
	"github.com/cepmap/otus-system-monitoring/internal/stats/loadavg"
)

type Collector struct {
	metrics   *metrics.Storage
	statTypes []pb.StatType
	avgPeriod time.Duration
}

func New(metrics *metrics.Storage, statTypes []pb.StatType, avgPeriod time.Duration) *Collector {
	return &Collector{
		metrics:   metrics,
		statTypes: statTypes,
		avgPeriod: avgPeriod,
	}
}

func (c *Collector) collectLoadAverage(timestamp time.Time) {
	if !config.DaemonConfig.Stats.LoadAverage {
		return
	}
	if stats, err := loadavg.GetStats(); err == nil {
		c.metrics.StoreLoadAverage(stats, timestamp)
	}
}

func (c *Collector) collectCPUStats(timestamp time.Time) {
	if !config.DaemonConfig.Stats.Cpu {
		return
	}
	if stats, err := cpu.GetCpuStat(); err == nil {
		c.metrics.StoreCPUStats(stats, timestamp)
	}
}

func (c *Collector) collectDisksLoad(timestamp time.Time) {
	if !config.DaemonConfig.Stats.DiskLoad {
		return
	}
	if stats, err := disksload.GetStats(); err == nil {
		c.metrics.StoreDisksLoad(stats, timestamp)
	}
}

func (c *Collector) collectDiskUsage(timestamp time.Time) {
	if !config.DaemonConfig.Stats.DiskInfo {
		return
	}
	if stats, err := diskstat.GetStats(); err == nil {
		c.metrics.StoreDiskUsage(stats, timestamp)
	}
}

func (c *Collector) CollectMetrics(timestamp time.Time) {
	for _, statType := range c.statTypes {
		switch statType {
		case pb.StatType_LOAD_AVERAGE:
			c.collectLoadAverage(timestamp)
		case pb.StatType_CPU_STATS:
			c.collectCPUStats(timestamp)
		case pb.StatType_DISKS_LOAD:
			c.collectDisksLoad(timestamp)
		case pb.StatType_DISK_USAGE:
			c.collectDiskUsage(timestamp)
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

func (c *Collector) prepareLoadAverageResponse(response *pb.StatsResponse) {
	if !config.DaemonConfig.Stats.LoadAverage {
		return
	}
	if avgStats := c.metrics.GetAverageLoadAverage(c.avgPeriod); avgStats != nil {
		response.LoadAverage = converter.LoadAverageToProto(avgStats)
	}
}

func (c *Collector) prepareCPUStatsResponse(response *pb.StatsResponse) {
	if !config.DaemonConfig.Stats.Cpu {
		return
	}
	if avgStats := c.metrics.GetAverageCPUStats(c.avgPeriod); avgStats != nil {
		response.CpuStats = converter.CPUStatToProto(avgStats)
	}
}

func (c *Collector) prepareDisksLoadResponse(response *pb.StatsResponse) {
	if !config.DaemonConfig.Stats.DiskLoad {
		return
	}
	if avgStats := c.metrics.GetAverageDisksLoad(c.avgPeriod); avgStats != nil {
		response.DisksLoad = converter.DisksLoadToProto(avgStats)
	}
}

func (c *Collector) prepareDiskUsageResponse(response *pb.StatsResponse) {
	if !config.DaemonConfig.Stats.DiskInfo {
		return
	}
	if stats := c.metrics.GetLatestDiskUsage(); stats != nil {
		response.DiskStats = converter.DiskStatsToProto(stats)
	}
}

func (c *Collector) PrepareResponse() *pb.StatsResponse {
	response := &pb.StatsResponse{
		Timestamp: time.Now().Unix(),
	}

	for _, statType := range c.statTypes {
		switch statType {
		case pb.StatType_LOAD_AVERAGE:
			c.prepareLoadAverageResponse(response)
		case pb.StatType_CPU_STATS:
			c.prepareCPUStatsResponse(response)
		case pb.StatType_DISKS_LOAD:
			c.prepareDisksLoadResponse(response)
		case pb.StatType_DISK_USAGE:
			c.prepareDiskUsageResponse(response)
		}
	}

	return response
}
