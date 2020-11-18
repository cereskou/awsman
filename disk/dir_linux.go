// +build linux

package disk

import "os"

//tempDir -
func tempDir() (dir string) {
	dir = "/tmp"
	return
}

//homeDir -
func homeDir() (dir string) {
	dir, _ = os.UserHomeDir()

	return
}
