package core

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/rs/zerolog/log"
)

var checks = []string{"down", "inactive", "not"}
var SECOND_REGEX, _ = regexp.Compile("(\\d+) seconds")

const MIN_SECOND = 4

func ClusterSanity(pkgName string, cmd string, clusterKeyword string) {
	clusters := tools.SearchCluster(clusterKeyword)
	for _, cluster := range clusters {
		log.Info().Str("Cluster", cluster).Msg("Sanity Check")
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
		log.Info().Str("Package", pkgVersion).Int("Count", count).Msg("Version Check")
	}
	if len(versionCountMap) != 1 {
		log.Warn().Str("Cluster", cluster).Msg("Multiple Versions Found")
	}
}

func VerifyStatus(cmd string, cluster string) {
	log.Info().Str("Command", cmd).Str("Cluster", cluster).Msg("Sanity Check")

	tools.NORMAL_PSSH.Run(cmd, cluster, 200, true)
	os.Chdir(config.OUTPUT_PATH)

	tools.PrintCommand("cat * | awk '{print $1,$2,$3}' | sort | uniq -c | sort -r")

	tools.RunIf("find  . -type f -empty | cut -c3-", func(output string) {
		if len(output) > 0 {
			log.Warn().Str("Cluster", cluster).Str("Output", output).Msg("Empty Files Found")
			tools.WriteClusterFile("empty", output)
		}
	})

	contentMap := util.ReadFileMap(config.OUTPUT_PATH, true)
	performBadStateChecks(contentMap)

	minUptime := getMinUptime(contentMap)
	if minUptime < MIN_SECOND {
		log.Warn().Str("Cluster", cluster).Int("Uptime", minUptime).Int("Threshold", MIN_SECOND).Msg("Probable Restart Detected")
	} else {
		log.Info().Str("Cluster", cluster).Int("Uptime", minUptime).Msg("Checks Complete")
	}

	if util.IsDebugMode() {
		VerifyNetworkParameters(cluster)
	}
}

func VerifyNetworkParameters(cluster string) {
	log.Info().Str("Cluster", cluster).Msg("Verifying Network Parameters")
	tools.Md5Checker("sudo sysctl -a | grep net | grep -v rss_key | grep -v nf_log", cluster)
}

/* Helpers */
func performBadStateChecks(contentMap map[string][]string) {
	for _, check := range checks {
		if keyWordLines, keyWordIps := extractKeywordLines(contentMap, check); len(keyWordLines) > 0 {
			log.Warn().Str("Check", check).Strs("Ips", keyWordIps).Strs("Lines", keyWordLines).Msg("Check Failed")
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
				log.Error().Str("Line", line).Err(err).Msg("Error Parsing Second")
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
