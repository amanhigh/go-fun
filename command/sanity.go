package commander

import (
	"fmt"
	"strings"
)

func VersionCheck(pkgNameCsv string, cluster string) {
	PrintBlue(fmt.Sprintf("Verifying Versions For Packages: %v on cluster %v", pkgNameCsv, cluster))
	packageList := strings.Split(pkgNameCsv, ",")

	cmd := fmt.Sprintf("dpkg -l | grep '%v'", strings.Join(packageList, `\|`))
	FastPssh.Run(cmd, cluster, DEFAULT_PARALELISM, true)

	for pkgVersion, count := range computeVersionCountMap() {
		PrintGreen(fmt.Sprintf("%v: %v", pkgVersion, count))
	}
}
func computeVersionCountMap() map[string]int {
	versionCountMap := map[string]int{}
	lines := ReadAllFiles(OUTPUT_PATH)
	for _, line := range lines {
		fields := strings.Fields(line)
		pkgName := fields[1]
		ver := fields[2]
		key := fmt.Sprintf("%v - %v", pkgName, ver)
		versionCountMap[key]++
	}
	return versionCountMap
}
