//go:build linux

package cpu

import (
	"github.com/cepmap/otus-system-monitoring/internal/models"
	tools "github.com/cepmap/otus-system-monitoring/internal/tools"
	"strings"
)

const (
	userPos   = 14
	systemPos = 16
	idlePos   = 19
)

func GetCpuStat() (*models.CPUStat, error) {
	res, err := tools.Exec("iostat", []string{"-c"})
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(res)

	user := tools.ParseFloat(fields[userPos])
	system := tools.ParseFloat(fields[systemPos])
	idle := tools.ParseFloat(fields[idlePos])

	return &models.CPUStat{
		User:   user,
		System: system,
		Idle:   idle,
	}, nil
}
