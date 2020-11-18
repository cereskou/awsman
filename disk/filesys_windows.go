// +build windows

package disk

import (
	"os"
	"syscall"
	"time"
)

//fileTime -
func fileTime(fi os.FileInfo) (createdate, updatedate, accessdate time.Time) {
	fs := fi.Sys().(*syscall.Win32FileAttributeData)
	createTime := time.Unix(0, fs.CreationTime.Nanoseconds())
	updateTime := time.Unix(0, fs.LastWriteTime.Nanoseconds())
	accessTime := time.Unix(0, fs.LastAccessTime.Nanoseconds())

	return createTime, updateTime, accessTime
}
