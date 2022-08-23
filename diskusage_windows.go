//go:build windows
// +build windows

package glog

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type DiskUsage struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// NewDiskUsages returns an object holding the disk usage of volumePath
// or nil in case of error (invalid path, etc)
func NewDiskUsage(volumePath string) (*DiskUsage, error) {
	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	h := syscall.MustLoadDLL("kernel32.dll")
	procGetDiskFreeSpaceExW := h.MustFindProc("GetDiskFreeSpaceExW")
	diskret, _, err := procGetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(volumePath))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	if diskret == 0 {
		return nil, err
	}
	return &DiskUsage{
		Path:        volumePath,
		Total:       uint64(lpTotalNumberOfBytes),
		Free:        uint64(lpTotalNumberOfFreeBytes),
		Used:        uint64(lpTotalNumberOfBytes) - uint64(lpTotalNumberOfFreeBytes),
		UsedPercent: (float64(lpTotalNumberOfBytes) - float64(lpTotalNumberOfFreeBytes)) / float64(lpTotalNumberOfBytes) * 100,
	}, nil
}
