// +build windows

package disk

import (
	"errors"
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows"
)

//diskUsage -
func diskUsage(path string) (free, total, avail uint64, err error) {
	if path == "" {
		path = "c:"
	} else {
		path = filepath.VolumeName(path)
	}

	var ptr *uint16
	ptr, err = windows.UTF16PtrFromString(path)
	if err != nil {
		return
	}
	err = windows.GetDiskFreeSpaceEx(ptr, &free, &total, &avail)
	if err != nil {
		if err == windows.ERROR_PATH_NOT_FOUND {
			err = errors.New("Path not found")
		} else {
			err = fmt.Errorf("An error has occurred. code %v", err)
		}
		return
	}
	return
}
