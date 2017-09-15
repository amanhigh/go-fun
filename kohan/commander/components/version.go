package components

import (
	"fmt"
	"time"
	"github.com/amanhigh/go-fun/util"
	"github.com/amanhigh/go-fun/kohan/commander"
	"github.com/amanhigh/go-fun/kohan/commander/tools"
)

func GetVersion(pkgName string, host string, versionType string) {
	switch versionType {
	case "dpkg":
		GetDpkgVersion(pkgName, host)
	case "latest":
		GetLatestVersion(pkgName, host)
	}
}

func GetDpkgVersion(pkgName string, host string) {
	util.PrintBlue(fmt.Sprintf("Fetching Config for %v from %v", pkgName, host))
	cmd := fmt.Sprintf(`ssh %v dpkg -l | grep "%v" | tail -1 | awk '{print $3}'`, host, pkgName)
	dpkgVersion := tools.RunCommandPrintError(cmd)
	versionString := fmt.Sprintf("\n%v - %v - HostVersion: %v", pkgName, dpkgVersion, util.FormatTime(time.Now(), util.PRINT_LAYOUT))
	util.PrintYellow(versionString)
	util.AppendFile(commander.RELEASE_FILE, versionString)
}

func GetLatestVersion(pkgName string, host string) {
	util.PrintBlue(fmt.Sprintf("Fetching LatestVersion for %v from %v", pkgName, host))
	cmd := fmt.Sprintf(`ssh %v "sudo apt-get update > /dev/null; apt-cache madison %v | head -1" | awk '{print $3}'`, host, pkgName)
	latestVersion := tools.RunCommandPrintError(cmd)
	versionString := fmt.Sprintf("\n%v - %v - LatestVersion ", pkgName, latestVersion)
	util.PrintYellow(versionString)
	util.AppendFile(commander.RELEASE_FILE, versionString)
}
