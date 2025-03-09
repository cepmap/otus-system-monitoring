package metrics

import (
	"math"

	"github.com/cepmap/otus-system-monitoring/internal/models"
)

func round(x float64) float64 {
	return math.Round(x*100) / 100
}

func averageLoadAverage(stats []*models.LoadAverage) *models.LoadAverage {
	if len(stats) == 0 {
		return nil
	}

	var sum models.LoadAverage
	for _, stat := range stats {
		sum.Load1Min += stat.Load1Min
		sum.Load5Min += stat.Load5Min
		sum.Load15Min += stat.Load15Min
	}

	count := float64(len(stats))
	return &models.LoadAverage{
		Load1Min:  round(sum.Load1Min / count),
		Load5Min:  round(sum.Load5Min / count),
		Load15Min: round(sum.Load15Min / count),
	}
}

func averageCPUStat(stats []*models.CPUStat) *models.CPUStat {
	if len(stats) == 0 {
		return nil
	}

	var sum models.CPUStat
	for _, stat := range stats {
		sum.User += stat.User
		sum.System += stat.System
		sum.Idle += stat.Idle
	}

	count := float64(len(stats))
	return &models.CPUStat{
		User:   round(sum.User / count),
		System: round(sum.System / count),
		Idle:   round(sum.Idle / count),
	}
}

func averageDisksLoad(stats []*models.DisksLoad) *models.DisksLoad {
	if len(stats) == 0 {
		return nil
	}

	diskSums := make(map[string]*struct {
		tpsSum float64
		kpsSum float64
		count  int
	})

	for _, stat := range stats {
		for _, disk := range stat.DisksLoad {
			if _, ok := diskSums[disk.FSName]; !ok {
				diskSums[disk.FSName] = &struct {
					tpsSum float64
					kpsSum float64
					count  int
				}{}
			}
			diskSums[disk.FSName].tpsSum += disk.Tps
			diskSums[disk.FSName].kpsSum += disk.Kps
			diskSums[disk.FSName].count++
		}
	}

	var result []models.DiskLoad
	for fsName, sums := range diskSums {
		result = append(result, models.DiskLoad{
			FSName: fsName,
			Tps:    round(sums.tpsSum / float64(sums.count)),
			Kps:    round(sums.kpsSum / float64(sums.count)),
		})
	}

	return &models.DisksLoad{DisksLoad: result}
}
