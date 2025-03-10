//go:build linux

package loadavg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStat(t *testing.T) {
	t.Run("test success get stats", func(t *testing.T) {
		load, err := GetStats()

		require.NoError(t, err)
		require.NotNil(t, load.Load1Min)
		require.IsType(t, float64(1), load.Load1Min)
		require.NotNil(t, load.Load5Min)
		require.IsType(t, float64(1), load.Load5Min)
		require.NotNil(t, load.Load1Min)
		require.IsType(t, float64(1), load.Load1Min)
	})
}
