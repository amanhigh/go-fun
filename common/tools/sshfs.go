package tools

import (
	"fmt"
	"github.com/amanhigh/go-fun/common/util"
	"os"
)

func SshfsMount(host string, remotePath string, localPath string) {
	os.MkdirAll(localPath, util.DIR_DEFAULT_PERM)
	SshfsUnmount(localPath)
	RunCommandPrintError(fmt.Sprintf("sshfs %v:%v %v", host, remotePath, localPath))
}

func SshfsUnmount(localPath string) {
	RunCommandPrintError(fmt.Sprintf("umount -f %v", localPath))
}
