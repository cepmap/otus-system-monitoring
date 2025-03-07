//go:build linux

package diskStat

import (
	"github.com/cepmap/otus-system-monitoring/internal/models"
	tools "github.com/cepmap/otus-system-monitoring/internal/tools"
	"strconv"
	"strings"
)

func parseUint(str string) (uint64, error) {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func GetDiskStats() (*models.DiskStats, error) {
	var output []models.DiskStat

	dfOut, err := getDiskInfo()
	if err != nil {
		return nil, err
	}

	dfInodeOut, err := getDiskInodeInfo()
	if err != nil {
		return nil, err
	}

	for i, disk := range dfOut {
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

func getDiskInfo() ([]string, error) {
	result, err := tools.Exec("df", []string{"-T", "-k", "--exclude-type=tmpfs", "--exclude-type=devtmpfs", "--exclude-type=udev"})
	if err != nil {
		return nil, err
	}
	return strings.Split(result, "\n")[1:], nil
}

func getDiskInodeInfo() ([]string, error) {
	result, err := tools.Exec("df", []string{"-T", "-k", "-i", "--exclude-type=tmpfs", "--exclude-type=devtmpfs", "--exclude-type=udev"})
	if err != nil {
		return nil, err
	}
	return strings.Split(result, "\n")[1:], nil
}
