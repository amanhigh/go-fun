package commander

import (
	"os"
)

func RecreateDir(path string,perm os.FileMode)  {
	os.RemoveAll(path)
	os.MkdirAll(path,perm)
}
