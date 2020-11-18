// +build windows

package disk

import "os"

//tempDir -
func tempDir() (dir string) {
	dir = os.Getenv("TEMP")

	return
}

//homeDir -
func homeDir() (dir string) {
	dir = os.Getenv("APPDATA")

	return
}
