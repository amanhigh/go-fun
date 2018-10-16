package tools

import (
	"fmt"

	"github.com/amanhigh/go-fun/util"
)

func Sync(srcHost string, srcDir string, targetDir string, targetHosts []string) {
	for _, target := range targetHosts {
		util.PrintSkyBlue(fmt.Sprintf("Syncing: %v:%v -> %v:%v", srcHost, srcDir, target, targetDir))
		RunCommandPrintError(fmt.Sprintf("ssh %v 'mkdir -p ~/sync'", target))
		RunCommandPrintError(fmt.Sprintf("ssh -A -t %v 'scp -o StrictHostKeyChecking=no -rC %v/* %v:~/sync'", srcHost, srcDir, target))
		RunCommandPrintError(fmt.Sprintf("ssh %v 'sudo mv ~/sync/* %v;rm -rf ~/sync'", target, targetDir))
	}
}

func SudoScp(fileName, srcDirectory, dstDirectory, dstHost string) {
	util.PrintSkyBlue(fmt.Sprintf("Moving File %v from %v to %v:%v", fileName, srcDirectory, dstHost, dstDirectory))
	RunCommandPrintError(fmt.Sprintf("scp -rC %v/%v %v:", srcDirectory, fileName, dstHost))
	RunCommandPrintError(fmt.Sprintf("ssh %v 'sudo mv ~/%v %v'", dstHost, fileName, dstDirectory))
}
