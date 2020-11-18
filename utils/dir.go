package utils

import (
	"os"
)

//MakeDirectory -
func MakeDirectory(dir string) error {
	//not exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
