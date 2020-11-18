// +build linux

package disk

import "syscall"

//diskUsage -
func diskUsage(path string) (free, total, avail uint64, err error) {
	if path == "" {
		path = "/"
	}

	var stat syscall.Statfs_t
	syscall.Statfs(path, &stat)

	free = stat.Bfree * uint64(stat.Bsize)
	total = stat.Blocks * uint64(stat.Bsize)
	avail = stat.Bavail * uint64(stat.Bsize)

	return
}
