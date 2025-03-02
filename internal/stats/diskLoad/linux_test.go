//go:build linux

package diskLoad

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStat(t *testing.T) {
	t.Run("test success get stats", func(t *testing.T) {
		disk, err := GetStats()

		require.NoError(t, err)
		require.NotNil(t, disk.Kps)
		require.IsType(t, float64(1), disk.Kps)
		require.NotNil(t, disk.Tps)
		require.IsType(t, float64(1), disk.Tps)
	})
}
