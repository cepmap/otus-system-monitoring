package converter

import (
	"testing"

	"github.com/cepmap/otus-system-monitoring/internal/models"
	"github.com/stretchr/testify/require"
)

func TestLoadAverageToProto(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		result := LoadAverageToProto(nil)
		require.Nil(t, result)
	})

	t.Run("valid input", func(t *testing.T) {
		input := &models.LoadAverage{
			Load1Min:  1.5,
			Load5Min:  2.0,
			Load15Min: 1.8,
		}
		result := LoadAverageToProto(input)
		require.NotNil(t, result)
		require.Equal(t, input.Load1Min, result.Load1Min)
		require.Equal(t, input.Load5Min, result.Load5Min)
		require.Equal(t, input.Load15Min, result.Load15Min)
	})
}

func TestCPUStatToProto(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		result := CPUStatToProto(nil)
		require.Nil(t, result)
	})

	t.Run("valid input", func(t *testing.T) {
		input := &models.CPUStat{
			User:   10.5,
			System: 5.2,
			Idle:   84.3,
		}
		result := CPUStatToProto(input)
		require.NotNil(t, result)
		require.Equal(t, input.User, result.User)
		require.Equal(t, input.System, result.System)
		require.Equal(t, input.Idle, result.Idle)
	})
}

func TestDisksLoadToProto(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		result := DisksLoadToProto(nil)
		require.Nil(t, result)
	})

	t.Run("valid input", func(t *testing.T) {
		input := &models.DisksLoad{
			DisksLoad: []models.DiskLoad{
				{
					FSName: "/dev/sda1",
					Tps:    100.5,
					Kps:    1024.0,
				},
				{
					FSName: "/dev/sdb1",
					Tps:    50.2,
					Kps:    512.0,
				},
			},
		}
		result := DisksLoadToProto(input)
		require.NotNil(t, result)
		require.Len(t, result.DisksLoad, len(input.DisksLoad))

		for i, disk := range input.DisksLoad {
			require.Equal(t, disk.FSName, result.DisksLoad[i].FsName)
			require.Equal(t, disk.Tps, result.DisksLoad[i].Tps)
			require.Equal(t, disk.Kps, result.DisksLoad[i].Kps)
		}
	})

	t.Run("empty disks list", func(t *testing.T) {
		input := &models.DisksLoad{
			DisksLoad: []models.DiskLoad{},
		}
		result := DisksLoadToProto(input)
		require.NotNil(t, result)
		require.Empty(t, result.DisksLoad)
	})
}

func TestDiskStatsToProto(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		result := DiskStatsToProto(nil)
		require.Nil(t, result)
	})

	t.Run("valid input", func(t *testing.T) {
		input := &models.DiskStats{
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
				{
					FileSystem: "/dev/sdb1",
					Usage: models.DiskUsage{
						Used:  1000000,
						Usage: "75%",
					},
					Inodes: models.InodeUsage{
						Used:  2000,
						Usage: "20%",
					},
				},
			},
		}
		result := DiskStatsToProto(input)
		require.NotNil(t, result)
		require.Len(t, result.DiskStats, len(input.DiskStats))

		for i, disk := range input.DiskStats {
			require.Equal(t, disk.FileSystem, result.DiskStats[i].Filesystem)
			require.Equal(t, disk.Usage.Used, result.DiskStats[i].Usage.Used)
			require.Equal(t, disk.Usage.Usage, result.DiskStats[i].Usage.Usage)
			require.Equal(t, disk.Inodes.Used, result.DiskStats[i].Inodes.Used)
			require.Equal(t, disk.Inodes.Usage, result.DiskStats[i].Inodes.Usage)
		}
	})

	t.Run("empty disks list", func(t *testing.T) {
		input := &models.DiskStats{
			DiskStats: []models.DiskStat{},
		}
		result := DiskStatsToProto(input)
		require.NotNil(t, result)
		require.Empty(t, result.DiskStats)
	})
}
