package tools

import "fmt"

func Sync(remoteHost string, srcDir string, targetDir string, targetHosts []string) {
	for _, target := range targetHosts {
		RunCommandPrintError(fmt.Sprintf("ssh %v 'mkdir -p ~/sync'", target))
		RunCommandPrintError(fmt.Sprintf("ssh -A -t %v 'scp -rC %v/* %v:~/sync'", remoteHost, srcDir, target))
		RunCommandPrintError(fmt.Sprintf("ssh %v 'sudo mv ~/sync/* %v;rm -rf ~/sync'", target, targetDir))
	}

}
