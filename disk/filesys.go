package disk

import (
	"os"
	"time"
)

//FileTime -
func FileTime(fi os.FileInfo) (createdate, updatedate, accessdate time.Time) {
	return fileTime(fi)
}
