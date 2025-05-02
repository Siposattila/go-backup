package disk

import "syscall"

type DiskUsage struct {
	stat *syscall.Statfs_t
}

func NewDiskUsage(volumePath string) *DiskUsage {
	var stat syscall.Statfs_t
	syscall.Statfs(volumePath, &stat)

	return &DiskUsage{&stat}
}

// Free returns total free bytes on file system
func (du *DiskUsage) Free() uint64 {
	return du.stat.Bfree * uint64(du.stat.Bsize)
}

// Available return total available bytes on file system to an unprivileged user
func (du *DiskUsage) Available() uint64 {
	return du.stat.Bavail * uint64(du.stat.Bsize)
}

// Size returns total size of the file system
func (du *DiskUsage) Size() uint64 {
	return uint64(du.stat.Blocks) * uint64(du.stat.Bsize)
}

// Used returns total bytes used in file system
func (du *DiskUsage) Used() uint64 {
	return du.Size() - du.Free()
}

// Usage returns percentage of use on the file system
func (du *DiskUsage) Usage() int {
	return int(du.Used()) * 100 / int(du.Size())
}
