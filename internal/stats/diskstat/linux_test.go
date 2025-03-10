//go:build linux

package diskstat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDiskStats(t *testing.T) {
	t.Run("test success get disk stats", func(t *testing.T) {
		stats, err := GetDiskStats()

		require.NoError(t, err)
		require.NotNil(t, stats)
		require.NotEmpty(t, stats.DiskStats)

		firstDisk := stats.DiskStats[0]

		require.NotEmpty(t, firstDisk.FileSystem)
		require.IsType(t, "", firstDisk.FileSystem)

		require.NotNil(t, firstDisk.Usage)
		require.IsType(t, uint64(1), firstDisk.Usage.Used)
		require.NotEmpty(t, firstDisk.Usage.Usage)
		require.IsType(t, "", firstDisk.Usage.Usage)

		require.NotNil(t, firstDisk.Inodes)
		require.IsType(t, uint64(1), firstDisk.Inodes.Used)
		require.NotEmpty(t, firstDisk.Inodes.Usage)
		require.IsType(t, "", firstDisk.Inodes.Usage)
	})
}

func TestParseUint(t *testing.T) {
	t.Run("test valid uint parsing", func(t *testing.T) {
		val, err := parseUint("123")
		require.NoError(t, err)
		require.Equal(t, uint64(123), val)
	})

	t.Run("test invalid uint parsing", func(t *testing.T) {
		val, err := parseUint("invalid")
		require.Error(t, err)
		require.Equal(t, uint64(0), val)
	})
}

func TestGetDiskInfo(t *testing.T) {
	t.Run("test get disk info success", func(t *testing.T) {
		info, err := getDiskInfo()
		require.NoError(t, err)
		require.NotEmpty(t, info)
		for _, line := range strings.Split(info, "\n") {
			fields := strings.Fields(line)
			require.GreaterOrEqual(t, len(fields), 6, "each line should have at least 6 fields")
		}
	})
}

func TestGetDiskInodeInfo(t *testing.T) {
	t.Run("test get disk inode info success", func(t *testing.T) {
		info, err := getDiskInodeInfo()
		require.NoError(t, err)
		require.NotEmpty(t, info)

		for _, line := range info {
			fields := strings.Fields(line)
			require.GreaterOrEqual(t, len(fields), 6, "each line should have at least 6 fields")
		}
	})
}
