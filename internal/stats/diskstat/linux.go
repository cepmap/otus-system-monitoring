//go:build linux

package diskstat

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cepmap/otus-system-monitoring/internal/models"
	tools "github.com/cepmap/otus-system-monitoring/internal/tools"
)

func parseUint(str string) (uint64, error) {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func GetDiskStats() (*models.DiskStats, error) {
	dfOut, err := getDiskInfo()
	if err != nil {
		return nil, fmt.Errorf("error getting disk info: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(dfOut), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected df output format: too few lines")
	}
	devices := lines[1:]
	output := make([]models.DiskStat, 0, len(devices))

	dfInodeOut, err := getDiskInodeInfo()
	if err != nil {
		return nil, err
	}

	for i, disk := range devices {
		diskArr := strings.Fields(disk)
		diskInodeArr := strings.Fields(dfInodeOut[i])

		if len(diskArr) < 6 || len(diskInodeArr) < 6 {
			continue
		}

		nDisk := models.DiskStat{
			FileSystem: diskArr[0],
			Usage:      models.DiskUsage{},
			Inodes:     models.InodeUsage{},
		}

		used, err := parseUint(diskArr[3])
		if err != nil {
			return nil, err
		}
		nDisk.Usage.Used = used
		nDisk.Usage.Usage = diskArr[5]

		inodesUsed, err := parseUint(diskInodeArr[3])
		if err != nil {
			return nil, err
		}
		nDisk.Inodes.Used = inodesUsed

		nDisk.Inodes.Usage = diskInodeArr[5]

		output = append(output, nDisk)
	}

	return &models.DiskStats{DiskStats: output}, nil
}

func getDiskInfo() (string, error) {
	result, err := tools.Exec("df", []string{
		"-T", "-k", "--exclude-type=tmpfs",
		"--exclude-type=devtmpfs", "--exclude-type=udev",
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

func getDiskInodeInfo() ([]string, error) {
	result, err := tools.Exec("df", []string{
		"-T", "-k", "-i", "--exclude-type=tmpfs",
		"--exclude-type=devtmpfs", "--exclude-type=udev",
	})
	if err != nil {
		return nil, err
	}
	lines := strings.Split(result, "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected df output format: too few lines")
	}
	return lines[1:], nil
}
