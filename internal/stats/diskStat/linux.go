//go:build linux

package diskStat

import (
	"github.com/cepmap/otus-system-monitoring/internal/models"
	tools "github.com/cepmap/otus-system-monitoring/internal/tools"
	"strconv"
	"strings"
)

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

		nDisk := models.DiskStat{}
		nDisk.Usage = models.DiskUsage{}
		nDisk.Inodes = models.InodeUsage{}

		val, err := strconv.ParseUint(diskArr[3], 10, 64)
		if err != nil {
			return nil, err
		}
		nDisk.Usage.Used = val
		nDisk.Usage.Usage = diskArr[5]

		val, err = strconv.ParseUint(diskInodeArr[3], 10, 64)
		if err != nil {
			return nil, err
		}
		nDisk.Inodes.Used = val
		val, err = strconv.ParseUint(diskInodeArr[3], 10, 64)
		if err != nil {
			return nil, err
		}
		nDisk.Inodes.Usage = diskInodeArr[5]
		output = append(output, nDisk)
	}

	return &models.DiskStats{output}, nil
}
func getDiskInfo() ([]string, error) {
	// Filesystem     Type  1K-blocks      Used Available Use% Mounted on
	result, err := tools.Exec("df", []string{"-T", "-k", "--exclude-type=tmpfs", "--exclude-type=devtmpfs", "--exclude-type=udev"})
	if err != nil {
		return nil, err
	}
	return strings.Split(result, "\n")[1:], nil
}

func getDiskInodeInfo() ([]string, error) {
	// Filesystem     Type   Inodes   IUsed   IFree IUse% Mounted on
	result, err := tools.Exec("df", []string{"-T", "-k", "-i", "--exclude-type=tmpfs", "--exclude-type=devtmpfs", "--exclude-type=udev"})
	if err != nil {
		return nil, err
	}
	return strings.Split(result, "\n")[1:], nil
}
