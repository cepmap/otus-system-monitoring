package converter

import (
	pb "github.com/cepmap/otus-system-monitoring/internal/api/stats_service"
	"github.com/cepmap/otus-system-monitoring/internal/models"
)

func LoadAverageToProto(la *models.LoadAverage) *pb.LoadAverage {
	if la == nil {
		return nil
	}
	return &pb.LoadAverage{
		Load1Min:  la.Load1Min,
		Load5Min:  la.Load5Min,
		Load15Min: la.Load15Min,
	}
}

func CPUStatToProto(cs *models.CPUStat) *pb.CPUStat {
	if cs == nil {
		return nil
	}
	return &pb.CPUStat{
		User:   cs.User,
		System: cs.System,
		Idle:   cs.Idle,
	}
}

func DisksLoadToProto(dl *models.DisksLoad) *pb.DisksLoad {
	if dl == nil {
		return nil
	}

	disks := make([]*pb.DiskLoad, len(dl.DisksLoad))
	for i, disk := range dl.DisksLoad {
		disks[i] = &pb.DiskLoad{
			FsName: disk.FSName,
			Tps:    disk.Tps,
			Kps:    disk.Kps,
		}
	}
	return &pb.DisksLoad{
		DisksLoad: disks,
	}
}

func DiskStatsToProto(ds *models.DiskStats) *pb.DiskStats {
	if ds == nil {
		return nil
	}

	diskStats := make([]*pb.DiskStat, len(ds.DiskStats))
	for i, diskStat := range ds.DiskStats {
		diskStats[i] = &pb.DiskStat{
			Filesystem: diskStat.FileSystem,
			Usage: &pb.DiskUsage{
				Used:  diskStat.Usage.Used,
				Usage: diskStat.Usage.Usage,
			},
			Inodes: &pb.InodeUsage{
				Used:  diskStat.Inodes.Used,
				Usage: diskStat.Inodes.Usage,
			},
		}
	}
	return &pb.DiskStats{
		DiskStats: diskStats,
	}
}
