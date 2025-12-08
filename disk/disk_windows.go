//go:build windows

package disk

import (
	"fmt"
	"syscall"
	"unsafe"
)

type DiskUsage struct {
	FreeBytes      uint64
	AvailableBytes uint64
	TotalBytes     uint64
	UsedBytes      uint64
	UsagePercent   int8
}

func NewDiskUsage(volumePath string) (*DiskUsage, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceExW := kernel32.NewProc("GetDiskFreeSpaceExW")

	lpFreeBytesAvailable := uint64(0)
	lpTotalNumberOfBytes := uint64(0)
	lpTotalNumberOfFreeBytes := uint64(0)

	volumePtr, _ := syscall.UTF16PtrFromString(volumePath)
	_, _, err := getDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(volumePtr)),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)),
	)
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("kernel32 getDiskFreeSpaceExW failed: %w", err)
	}

	used := lpTotalNumberOfBytes - lpTotalNumberOfFreeBytes

	return &DiskUsage{
		FreeBytes:      lpTotalNumberOfFreeBytes,
		AvailableBytes: lpFreeBytesAvailable,
		TotalBytes:     lpTotalNumberOfBytes,
		UsedBytes:      used,
		UsagePercent:   int8(used * 100 / lpTotalNumberOfBytes),
	}, nil
}
