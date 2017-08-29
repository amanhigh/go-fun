package commander

import (
	"fmt"
	"io/ioutil"
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
	versionString := fmt.Sprintf("%v - %v - HostVersion: %v", pkgName, dpkgVersion, time.Now().Format(LAYOUT))
	PrintYellow(versionString)
	ioutil.WriteFile(RELEASE_FILE, []byte(versionString), DEFAULT_PERM)
}

func GetLatestVersion(pkgName string, host string) {
	PrintBlue(fmt.Sprintf("Fetching LatestVersion for %v from %v", pkgName, host))
	cmd := fmt.Sprintf(`ssh %v "sudo apt-get update > /dev/null; apt-cache madison %v | head -1" | awk '{print $3}'`, host, pkgName)
	latestVersion := RunCommandPrintError(cmd)
	versionString := fmt.Sprintf("%v - %v - LatestVersion: %v", pkgName, latestVersion, time.Now().Format(LAYOUT))
	PrintYellow(versionString)
	ioutil.WriteFile(RELEASE_FILE, []byte(versionString), DEFAULT_PERM)
}
