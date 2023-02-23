package files

import (
	"os"
)

func CreateDirectoryIfNotExists(p string, perm os.FileMode) error {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return os.Mkdir(p, perm)
	}

	return nil
}
