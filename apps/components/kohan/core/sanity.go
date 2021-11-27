package core

import (
	"fmt"
	"github.com/amanhigh/go-fun/apps/models/config"
	"github.com/fatih/color"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/apps/common/tools"
	"github.com/amanhigh/go-fun/apps/common/util"
	log "github.com/sirupsen/logrus"
)

var checks = []string{"down", "inactive", "not"}
var SECOND_REGEX, _ = regexp.Compile("(\\d+) seconds")

const MIN_SECOND = 4

func ClusterSanity(pkgName string, cmd string, clusterKeyword string) {
	clusters := tools.SearchCluster(clusterKeyword)
	for _, cluster := range clusters {
		color.Yellow("Processing: " + cluster)
		if cmd != "" {
			VerifyStatus(cmd, cluster)
		}
		VersionCheck(pkgName, cluster)
	}
}

func VersionCheck(pkgNameCsv string, cluster string) {
	packageList := strings.Split(pkgNameCsv, ",")

	cmd := fmt.Sprintf("dpkg -l | grep '%v'", strings.Join(packageList, `\|`))
	tools.FastPssh.Run(cmd, cluster, config.DEFAULT_PARALELISM, true)

	versionCountMap := computeVersionCountMap()
	for pkgVersion, count := range versionCountMap {
		color.Green("%v : %v", pkgVersion, count)
	}
	if len(versionCountMap) != 1 {
		color.Red("Multiple Versions Found on %v", cluster)
	}
}

func VerifyStatus(cmd string, cluster string) {
	color.Blue("Running Sanity on Cluster: " + cluster)

	tools.NORMAL_PSSH.Run(cmd, cluster, 200, true)
	os.Chdir(config.OUTPUT_PATH)

	tools.PrintCommand("cat * | awk '{print $1,$2,$3}' | sort | uniq -c | sort -r")

	tools.RunIf("find  . -type f -empty | cut -c3-", func(output string) {
		if len(output) > 0 {
			color.Red(fmt.Sprintf("Empty Files Found:\n%v", output))
			tools.WriteClusterFile("empty", output)
		}
	})

	contentMap := util.ReadFileMap(config.OUTPUT_PATH, true)
	performBadStateChecks(contentMap)

	minUptime := getMinUptime(contentMap)
	if minUptime < MIN_SECOND {
		color.Red("Probable Restart Detected. Second: %v", minUptime)
	} else {
		color.Green("Checks Complete, Min Uptime (seconds): %v", minUptime)
	}

	//TODO:Move out of Debug Mode.
	if util.IsDebugMode() {
		VerifyNetworkParameters(cluster)
	}
}

func VerifyNetworkParameters(cluster string) {
	color.Yellow("\nVerifying Network Parameters. Cluster: " + cluster)
	tools.Md5Checker("sudo sysctl -a | grep net | grep -v rss_key | grep -v nf_log", cluster)
}

/* Helpers */
func performBadStateChecks(contentMap map[string][]string) {
	for _, check := range checks {
		if keyWordLines, keyWordIps := extractKeywordLines(contentMap, check); len(keyWordLines) > 0 {
			color.Blue("Check Failed: " + check)
			color.Red(strings.Join(keyWordLines, "\n"))
			tools.WriteClusterFile(check, strings.Join(keyWordIps, "\n"))
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
	lines := util.ReadAllFiles(config.OUTPUT_PATH)
	for _, line := range lines {
		fields := strings.Fields(line)
		pkgName := fields[1]
		ver := fields[2]
		key := fmt.Sprintf("%v - %v", pkgName, ver)
		versionCountMap[key]++
	}
	return versionCountMap
}
