package loadAvg

import (
	"strings"

	"github.com/cepmap/otus-system-monitoring/internal/models"
	tools "github.com/cepmap/otus-system-monitoring/internal/tools"
)

const (
	load1MinPos  = 0
	load5MinPos  = 1
	load15MinPos = 2
)

func GetStatsOs() (*models.LoadAverage, error) {
	res, err := tools.Exec("cat", []string{"/proc/loadavg"})
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(res)

	load1Min := tools.ParseFloat(fields[load1MinPos])
	load5Min := tools.ParseFloat(fields[load5MinPos])
	load15Min := tools.ParseFloat(fields[load15MinPos])

	return &models.LoadAverage{
		Load1Min:  load1Min,
		Load5Min:  load5Min,
		Load15Min: load15Min,
	}, nil
}
