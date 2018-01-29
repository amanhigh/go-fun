package components

import (
	"fmt"
	"github.com/Flipkart/elb/scripts/kohan/util"
	log "github.com/Sirupsen/logrus"
	. "github.com/amanhigh/go-fun/kohan/commander/tools"
	. "github.com/amanhigh/go-fun/util"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var checks = []string{"down", "inactive", "not"}
var SECOND_REGEX, _ = regexp.Compile("(\\d+) seconds")

const MIN_SECOND = 4

func ClusterSanity(pkgName string, cluster string) {
	VerifyStatus(pkgName, cluster)
	VersionCheck(pkgName, cluster)
}

func VersionCheck(pkgNameCsv string, cluster string) {
	PrintBlue(fmt.Sprintf("Verifying Versions For Packages: %v on cluster %v", pkgNameCsv, cluster))
	packageList := strings.Split(pkgNameCsv, ",")

	cmd := fmt.Sprintf("dpkg -l | grep '%v'", strings.Join(packageList, `\|`))
	FastPssh.Run(cmd, cluster, DEFAULT_PARALELISM, true)

	for pkgVersion, count := range computeVersionCountMap() {
		PrintGreen(fmt.Sprintf("%v: %v", pkgVersion, count))
	}
}

func VerifyStatus(pkgName string, cluster string) {
	PrintBlue("Running Sanity on Cluster: " + cluster)

	//pr $cluster 100 "$cmd;sudo /etc/init.d/nsca status;sudo /etc/init.d/cosmos-jmx status" 10 > /dev/null;
	cmd := util.KohanConfig.PackageCheckCommands[pkgName]
	NORMAL_PSSH.Run(cmd, cluster, 200, true)
	os.Chdir(OUTPUT_PATH)

	PrintCommand("cat * | awk '{print $1,$2,$3}' | sort | uniq -c | sort -r")

	RunIf("find  . -type f -empty | cut -c3-", func(output string) {
		if len(output) > 0 {
			PrintRed(fmt.Sprintf("Empty Files Found:\n%v", output))
			WriteClusterFile("empty", output)
		}
	})

	contentMap := ReadFileMap(OUTPUT_PATH)
	performBadStateChecks(contentMap)

	minFound := performSecondsCheck(contentMap)
	PrintBlue(fmt.Sprintf("Second Check Complete. Min Second Detected: %v", minFound))

	//TODO:Move out of Debug Mode.
	if IsDebugMode() {
		VerifyNetworkParameters(cluster)
	}
}

func VerifyNetworkParameters(cluster string) {
	PrintYellow("\nVerifying Network Parameters. Cluster: " + cluster)
	Md5Checker("sudo sysctl -a | grep net | grep -v rss_key | grep -v nf_log", cluster)
}

/* Helpers */
func performBadStateChecks(contentMap map[string][]string) {
	for _, check := range checks {
		if keyWordLines, keyWordIps := extractKeywordLines(contentMap, check); len(keyWordLines) > 0 {
			PrintBlue("Check Failed: " + check)
			PrintRed(strings.Join(keyWordLines, "\n"))
			WriteClusterFile(check, strings.Join(keyWordIps, "\n"))
		}
	}
	PrintGreen(fmt.Sprintf("Checks Complete: %v", checks))
}

func performSecondsCheck(contentMap map[string][]string) int {
	minFound := math.MaxInt64
	if lines, _ := extractKeywordLines(contentMap, "seconds"); len(lines) > 0 {
		for _, line := range lines {
			matchString := SECOND_REGEX.FindStringSubmatch(line)
			if second, err := strconv.Atoi(matchString[1]); err == nil {
				if second < minFound {
					minFound = second
				}
			} else {
				log.WithFields(log.Fields{"Error": err}).Error("Error Parsing Second")
			}
		}

		if minFound < MIN_SECOND {
			PrintRed(fmt.Sprintf("Probable Restart Detected. Second: %v", minFound))
		}
	}
	return minFound
}
func extractKeywordLines(contentMap map[string][]string, keyWord string) ([]string, []string) {
	keyWordLines := []string{}
	keyWordIps := []string{}
	for ip, lines := range contentMap {
		for _, line := range lines {
			if ok := strings.Contains(line, keyWord); ok {
				keyWordLines = append(keyWordLines, line)
				keyWordIps = append(keyWordIps, ip)
			}
		}
	}
	return keyWordLines, keyWordIps
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
