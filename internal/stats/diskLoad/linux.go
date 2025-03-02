//go:build linux

package diskLoad

import (
	"github.com/cepmap/otus-system-monitoring/internal/models"
	tools "github.com/cepmap/otus-system-monitoring/internal/tools"
	"strings"
)

const (
	kbtPos = 29
	tpsPos = 30
)

func GetDiskLoad() (*models.DiskLoad, error) {
	res, err := tools.Exec("iostat", []string{""})
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(res)

	kbt := tools.ParseFloat(fields[kbtPos])
	tps := tools.ParseFloat(fields[tpsPos])

	return &models.DiskLoad{
		Kps: kbt,
		Tps: tps,
	}, nil
}
