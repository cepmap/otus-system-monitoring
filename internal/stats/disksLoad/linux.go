//go:build linux

package disksLoad

import (
	"strings"

	"github.com/cepmap/otus-system-monitoring/internal/models"
	tools "github.com/cepmap/otus-system-monitoring/internal/tools"
)

func GetDisksLoad() (*models.DisksLoad, error) {
	var disksLoad []models.DiskLoad

	res, err := tools.ExecCommand("iostat", []string{""})
	if err != nil {
		return nil, err
	}
	lines := strings.Split(res, "\n")
	for i := range lines {
		fields := strings.Fields(lines[i])
		if len(fields) < 8 || fields[0] == "Device" {
			continue
		}
		nDiskLoad := models.DiskLoad{
			FSName: fields[0],
			Kps:    tools.ParseFloat(fields[2]),
			Tps:    tools.ParseFloat(fields[1]),
		}
		disksLoad = append(disksLoad, nDiskLoad)
	}
	return &models.DisksLoad{DisksLoad: disksLoad}, nil
}
