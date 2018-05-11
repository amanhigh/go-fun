package tools

import (
	"fmt"

	"github.com/amanhigh/go-fun/util"
)

func Sync(srcHost string, srcDir string, targetDir string, targetHosts []string) {
	for _, target := range targetHosts {
		util.PrintSkyBlue(fmt.Sprintf("Syncing: %v:%v -> %v:%v", srcHost, srcDir, target, targetDir))
		RunCommandPrintError(fmt.Sprintf("ssh %v 'mkdir -p ~/sync'", target))
		RunCommandPrintError(fmt.Sprintf("ssh -A -t %v 'scp -rC %v/* %v:~/sync'", srcHost, srcDir, target))
		RunCommandPrintError(fmt.Sprintf("ssh %v 'sudo mv ~/sync/* %v;rm -rf ~/sync'", target, targetDir))
	}

}
