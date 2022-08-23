//go:build freebsd || linux || darwin
// +build freebsd linux darwin

package glog

import (
	"strconv"

	"golang.org/x/sys/unix"
)

// DiskUsage contains usage data and provides user-friendly access methods
type DiskUsage struct {
	Path              string  `json:"path"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

// NewDiskUsages returns an object holding the disk usage of volumePath
// or nil in case of error (invalid path, etc)
func NewDiskUsage(volumePath string) (*DiskUsage, error) {
	stat := unix.Statfs_t{}
	err := unix.Statfs(volumePath, &stat)
	if err != nil {
		return nil, err
	}
	bsize := uint64(stat.Bsize)

	ret := &DiskUsage{
		Path:        unescapeFstab(volumePath),
		Total:       stat.Blocks * bsize,
		Free:        stat.Bavail * bsize,
		InodesTotal: stat.Files,
		InodesFree:  stat.Ffree,
	}

	// if could not get InodesTotal, return empty
	if ret.InodesTotal < ret.InodesFree {
		return ret, nil
	}

	ret.InodesUsed = (ret.InodesTotal - ret.InodesFree)
	ret.Used = (uint64(stat.Blocks) - uint64(stat.Bfree)) * uint64(bsize)

	if ret.InodesTotal == 0 {
		ret.InodesUsedPercent = 0
	} else {
		ret.InodesUsedPercent = (float64(ret.InodesUsed) / float64(ret.InodesTotal)) * 100.0
	}

	if (ret.Used + ret.Free) == 0 {
		ret.UsedPercent = 0
	} else {
		// We don't use ret.Total to calculate percent.
		// see https://github.com/shirou/gopsutil/issues/562
		ret.UsedPercent = (float64(ret.Used) / float64(ret.Used+ret.Free)) * 100.0
	}

	return ret, nil
}

// Unescape escaped octal chars (like space 040, ampersand 046 and backslash 134) to their real value in fstab fields issue#555
func unescapeFstab(path string) string {
	escaped, err := strconv.Unquote(`"` + path + `"`)
	if err != nil {
		return path
	}
	return escaped
}
