package commander

import (
	"fmt"
	"io/ioutil"
	"time"
)

func GetDpkgVersion(pkgName string, host string) {
	PrintBlue(fmt.Sprintf("Fetching Config for %v from %v", pkgName, host))
	cmd := fmt.Sprintf(`ssh %v dpkg -l | grep "%v" | tail -1 | awk '{print $3}'`, host, pkgName)
	dpkgVersion := RunCommandPrintError(cmd)
	versionString := fmt.Sprintf("%v - %v - HostVersion: %v", pkgName, dpkgVersion, time.Now().Format(LAYOUT))
	PrintYellow(versionString)
	ioutil.WriteFile(RELEASE_FILE, []byte(versionString), DEFAULT_PERM)
}
