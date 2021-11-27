package core

import (
	"fmt"
	util2 "github.com/amanhigh/go-fun/apps/common/util"
	"github.com/amanhigh/go-fun/apps/models/config"
	"time"

	"github.com/amanhigh/go-fun/apps/common/tools"
	"github.com/amanhigh/go-fun/util"
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
	util.PrintBlue(fmt.Sprintf("Fetching Config for %v from %v", pkgName, host))
	cmd := fmt.Sprintf(`ssh %v dpkg -l | grep "%v" | tail -1 | awk '{print $3}'`, host, pkgName)
	dpkgVersion := tools.RunCommandPrintError(cmd)
	versionString := fmt.Sprintf("\n%v - %v - HostVersion: %v", pkgName, dpkgVersion, util2.FormatTime(time.Now(), util2.PRINT_LAYOUT))
	util.PrintYellow(versionString)
	util2.AppendFile(config.RELEASE_FILE, versionString)
}

func GetLatestVersion(pkgName string, host string, comment string) {
	util.PrintBlue(fmt.Sprintf("Fetching LatestVersion for %v from %v", pkgName, host))
	cmd := fmt.Sprintf(`ssh %v "sudo apt-get update > /dev/null; apt-cache madison %v | head -1" | awk '{print $3}'`, host, pkgName)
	latestVersion := tools.RunCommandPrintError(cmd)
	versionString := fmt.Sprintf("\n%v - %v - LatestVersion [ %v ]", pkgName, latestVersion, comment)
	util.PrintYellow(versionString)
	util2.AppendFile(config.RELEASE_FILE, versionString)
}
