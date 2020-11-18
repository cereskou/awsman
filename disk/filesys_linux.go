// +build linux

package disk

import (
	"os"
	"syscall"
	"time"
)

//fileTime -
func fileTime(fi os.FileInfo) (createdate, updatedate, accessdate time.Time) {
	fs := fi.Sys().(*syscall.Stat_t)
	createTime := time.Unix(0, fs.Ctim)
	updateTime := time.Unix(0, fs.Mtim)
	accessTime := time.Unix(0, fs.Atim)

	return createTime, updateTime, accessTime
}
