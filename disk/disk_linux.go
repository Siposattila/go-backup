//go:build linux

package disk

import (
	"fmt"
	"syscall"
)

type DiskUsage struct {
	FreeBytes      uint64
	AvailableBytes uint64
	TotalBytes     uint64
	UsedBytes      uint64
	UsagePercent   int8
}

func NewDiskUsage(volumePath string) (*DiskUsage, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(volumePath, &stat); err != nil {
		return nil, fmt.Errorf("statfs syscall failed: %w", err)
	}

	return &DiskUsage{
		FreeBytes:      stat.Bfree * uint64(stat.Bsize),
		AvailableBytes: stat.Bavail * uint64(stat.Bsize),
		TotalBytes:     uint64(stat.Blocks) * uint64(stat.Bsize),
		UsedBytes:      (uint64(stat.Blocks) * uint64(stat.Bsize)) - (stat.Bfree * uint64(stat.Bsize)),
		UsagePercent:   int8(int32((uint64(stat.Blocks)*uint64(stat.Bsize))-(stat.Bfree*uint64(stat.Bsize))) * 100 / int32(uint64(stat.Blocks)*uint64(stat.Bsize))),
	}, nil
}
