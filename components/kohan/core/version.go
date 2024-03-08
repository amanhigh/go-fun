package core

import (
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/fatih/color"
)

func GetVersion(pkgName string, host string, versionType string, comment string) {
	switch versionType {
	case "dpkg":
		GetDpkgVersion(pkgName, host)
	case "latest":
		GetLatestVersion(pkgName, host, comment)
	}
}

func GetDpkgVersion(pkgName string, host string) {
	color.Blue("Fetching Config for %v from %v", pkgName, host)
	cmd := fmt.Sprintf(`ssh %v dpkg -l | grep "%v" | tail -1 | awk '{print $3}'`, host, pkgName)
	dpkgVersion := tools.RunCommandPrintError(cmd)
	versionString := fmt.Sprintf("\n%v - %v - HostVersion: %v", pkgName, dpkgVersion, util.FormatTime(time.Now(), util.PRINT_LAYOUT))
	color.Yellow(versionString)
	util.AppendFile(config.RELEASE_FILE, versionString)
}

func GetLatestVersion(pkgName string, host string, comment string) {
	color.Blue("Fetching LatestVersion for %v from %v", pkgName, host)
	cmd := fmt.Sprintf(`ssh %v "sudo apt-get update > /dev/null; apt-cache madison %v | head -1" | awk '{print $3}'`, host, pkgName)
	latestVersion := tools.RunCommandPrintError(cmd)
	versionString := fmt.Sprintf("\n%v - %v - LatestVersion [ %v ]", pkgName, latestVersion, comment)
	color.Yellow(versionString)
	util.AppendFile(config.RELEASE_FILE, versionString)
}
