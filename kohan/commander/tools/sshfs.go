package tools

import (
	"fmt"
	"os"

	"github.com/amanhigh/go-fun/util"
)

func SshfsMount(host string, remotePath string, localPath string) {
	os.MkdirAll(localPath, util.DIR_DEFAULT_PERM)
	SshfsUnmount(localPath)
	RunCommandPrintError(fmt.Sprintf("sshfs %v:%v %v", host, remotePath, localPath))
}

func SshfsUnmount(localPath string) {
	RunCommandPrintError(fmt.Sprintf("umount %v", localPath))
}
