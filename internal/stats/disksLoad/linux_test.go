//go:build linux

package disksLoad

import (
	"testing"

	"github.com/cepmap/otus-system-monitoring/internal/models"
	"github.com/cepmap/otus-system-monitoring/internal/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTools struct {
	mock.Mock
}

func (m *MockTools) Exec(command string, args []string) (string, error) {
	return `Device            tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
sda              1.23    456.78       789.01          0.00     123456     789012          0
sdb              4.56    789.01       123.45          0.00     456789     123456          0`, nil
}

func TestGetDisksLoad(t *testing.T) {
	t.Run("should parse iostat output correctly", func(t *testing.T) {
		mockTools := new(MockTools)
		tools.ExecCommand = mockTools.Exec
		mockTools.On("Exec", "iostat", []string{""}).Return([]string{""})

		result, err := GetDisksLoad()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.DisksLoad, 2)

		assert.Equal(t, "sda", result.DisksLoad[0].FSName)
		assert.Equal(t, 1.23, result.DisksLoad[0].Tps)
		assert.Equal(t, 456.78, result.DisksLoad[0].Kps)

		assert.Equal(t, "sdb", result.DisksLoad[1].FSName)
		assert.Equal(t, 4.56, result.DisksLoad[1].Tps)
		assert.Equal(t, 789.01, result.DisksLoad[1].Kps)
	})
}

func TestGetStats(t *testing.T) {
	t.Run("should return disk stats", func(t *testing.T) {
		mockTools := new(MockTools)
		tools.ExecCommand = mockTools.Exec
		mockTools.On("Exec", "iostat", []string{""}).Return([]string{""})

		result, err := GetStats()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.IsType(t, &models.DisksLoad{}, result)
	})
}
