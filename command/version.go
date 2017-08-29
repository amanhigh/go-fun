package commander

import (
	"fmt"
	"time"
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
	PrintBlue(fmt.Sprintf("Fetching Config for %v from %v", pkgName, host))
	cmd := fmt.Sprintf(`ssh %v dpkg -l | grep "%v" | tail -1 | awk '{print $3}'`, host, pkgName)
	dpkgVersion := RunCommandPrintError(cmd)
	versionString := fmt.Sprintf("\n%v - %v - HostVersion: %v", pkgName, dpkgVersion, time.Now().Format(LAYOUT))
	PrintYellow(versionString)
	AppendFile(RELEASE_FILE,versionString)
}

func GetLatestVersion(pkgName string, host string) {
	PrintBlue(fmt.Sprintf("Fetching LatestVersion for %v from %v", pkgName, host))
	cmd := fmt.Sprintf(`ssh %v "sudo apt-get update > /dev/null; apt-cache madison %v | head -1" | awk '{print $3}'`, host, pkgName)
	latestVersion := RunCommandPrintError(cmd)
	versionString := fmt.Sprintf("\n%v - %v - LatestVersion ", pkgName, latestVersion)
	PrintYellow(versionString)
	AppendFile(RELEASE_FILE,versionString)
}
