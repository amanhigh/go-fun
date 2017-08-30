package commander

import (
	"fmt"
	"strings"
	"os"
)

var checks = []string{"down", "inactive"}

func VersionCheck(pkgNameCsv string, cluster string) {
	PrintBlue(fmt.Sprintf("Verifying Versions For Packages: %v on cluster %v", pkgNameCsv, cluster))
	packageList := strings.Split(pkgNameCsv, ",")

	cmd := fmt.Sprintf("dpkg -l | grep '%v'", strings.Join(packageList, `\|`))
	FastPssh.Run(cmd, cluster, DEFAULT_PARALELISM, true)

	for pkgVersion, count := range computeVersionCountMap() {
		PrintGreen(fmt.Sprintf("%v: %v", pkgVersion, count))
	}
}

func VerifyStatus(cmd string, cluster string) {
	PrintBlue("Running on Cluster: " + cluster)

	//pr $cluster 100 "$cmd;sudo /etc/init.d/nsca status;sudo /etc/init.d/cosmos-jmx status" 10 > /dev/null;
	FastPssh.Run(cmd, cluster, 200, true)
	os.Chdir(OUTPUT_PATH)

	PrintBlue("Summary")
	PrintCommand("cat * | awk '{print $1,$2,$3}' | sort | uniq -c | sort -r")

	RunIf("find  . -type f -empty | cut -c3-", func(output string) {
		PrintRed(fmt.Sprintf("Empty Files Found:\n%v", output))
		WriteClusterFile("empty",output)
	})

	//	sc down;
	//	sc inactive;
	//		echo -en "\033[1;34m Extracting Bad States \033[0m \n"
	//	sc out | awk -F: '{print $1}' | cut -c 3- > $cluster_path/oor.txt;
	//	sc not | awk -F: '{print $1}' | cut -c 3- > $cluster_path/not.txt;
	//	sc inactive | awk -F: '{print $1}' | cut -c 3- > $cluster_path/inactive.txt;

	//echo -en "\033[1;34m Seconds \033[0m \n"
	//grep -inrR "seconds" . 2> /dev/null | head -2
	//

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
