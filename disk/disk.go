package disk

import (
	"syscall"

	"github.com/Siposattila/go-backup/log"
)

type DiskUsage struct {
	stat *syscall.Statfs_t
}

func NewDiskUsage(volumePath string) *DiskUsage {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(volumePath, &stat); err != nil {
		log.GetLogger().Fatal("Statfs syscall failed: ", err.Error())
	}

	return &DiskUsage{&stat}
}

func (du *DiskUsage) Free() uint64 {
	return du.stat.Bfree * uint64(du.stat.Bsize)
}

func (du *DiskUsage) Available() uint64 {
	return du.stat.Bavail * uint64(du.stat.Bsize)
}

func (du *DiskUsage) Size() uint64 {
	return uint64(du.stat.Blocks) * uint64(du.stat.Bsize)
}

func (du *DiskUsage) Used() uint64 {
	return du.Size() - du.Free()
}

func (du *DiskUsage) Usage() int32 {
	return int32(du.Used()) * 100 / int32(du.Size())
}
