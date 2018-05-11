package components

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	. "github.com/amanhigh/go-fun/kohan/commander/tools"
	. "github.com/amanhigh/go-fun/util"
)

var checks = []string{"down", "inactive", "not"}
var SECOND_REGEX, _ = regexp.Compile("(\\d+) seconds")

const MIN_SECOND = 4

func ClusterSanity(pkgName string, cmd string, cluster string) {
	if cmd != "" {
		VerifyStatus(cmd, cluster)
	}
	VersionCheck(pkgName, cluster)
}

func VersionCheck(pkgNameCsv string, cluster string) {
	packageList := strings.Split(pkgNameCsv, ",")

	cmd := fmt.Sprintf("dpkg -l | grep '%v'", strings.Join(packageList, `\|`))
	FastPssh.Run(cmd, cluster, DEFAULT_PARALELISM, true)

	versionCountMap := computeVersionCountMap()
	for pkgVersion, count := range versionCountMap {
		PrintGreen(fmt.Sprintf("%v : %v", pkgVersion, count))
	}
	if len(versionCountMap) != 1 {
		PrintRed(fmt.Sprintf("Multiple Versions Found on %v", cluster))
	}
}

func VerifyStatus(cmd string, cluster string) {
	PrintBlue("Running Sanity on Cluster: " + cluster)

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

	minUptime := getMinUptime(contentMap)
	if minUptime < MIN_SECOND {
		PrintRed(fmt.Sprintf("Probable Restart Detected. Second: %v", minUptime))
	} else {
		PrintGreen(fmt.Sprintf("Checks Complete, Min Uptime (seconds): %v", minUptime))
	}

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
}

func getMinUptime(contentMap map[string][]string) int {
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
