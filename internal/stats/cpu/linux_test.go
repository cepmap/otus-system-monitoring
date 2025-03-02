//go:build linux

package cpu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStat(t *testing.T) {
	t.Run("test success get stats", func(t *testing.T) {
		cpu, err := GetStats()

		require.NoError(t, err)
		require.NotNil(t, cpu.User)
		require.IsType(t, float64(1), cpu.User)
		require.NotNil(t, cpu.System)
		require.IsType(t, float64(1), cpu.System)
		require.NotNil(t, cpu.Idle)
		require.IsType(t, float64(1), cpu.Idle)
	})
}
