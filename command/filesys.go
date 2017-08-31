package commander

import (
	"os"
)

func RecreateDir(path string) {
	os.RemoveAll(path)
	os.MkdirAll(path, DIR_DEFAULT_PERM)
}
