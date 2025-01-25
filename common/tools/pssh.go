package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

const (
	defaultTimeout = 200
)

var FastPssh = Pssh{20, config.OUTPUT_PATH, config.ERROR_PATH, false}
var NORMAL_PSSH = Pssh{30, config.OUTPUT_PATH, config.ERROR_PATH, false}
var DisplayPssh = Pssh{10, config.OUTPUT_PATH, config.ERROR_PATH, true}
var SlowPssh = Pssh{240, config.OUTPUT_PATH, config.ERROR_PATH, false}

type Pssh struct {
	Timeout       int
	outputPath    string
	errorPath     string
	displayOutput bool
}

func (p *Pssh) Run(cmd string, cluster string, parallelism int, disableOutput bool) {
	clearOutputPaths()

	psshCmd := fmt.Sprintf(`script %v pssh -h %v -t %v -o %v -e %v %v -p %v '%v'`,
		config.CONSOLE_FILE, getClusterFile(cluster), p.Timeout, p.outputPath, p.errorPath, p.getDisplayFlag(), parallelism, cmd)
	if disableOutput {
		RunCommandPrintError(psshCmd)
	} else {
		log.Info().Str("Cluster", cluster).Int("Parallelism", parallelism).Str("CMD", psshCmd).Msg("Running Parallel SSH")
		LiveCommand(psshCmd)
	}

	RunIf(fmt.Sprintf("grep FAILURE %v", getClusterFile("console")), func(output string) {
		PrintCommand(fmt.Sprintf("grep SUCCESS %v | awk '{print $4}' > %v", getClusterFile("console"), getClusterFile("pass")))
		PrintCommand(fmt.Sprintf("grep FAILURE %v | awk '{print $4}' > %v", getClusterFile("console"), getClusterFile("fail")))
		log.Warn().Str("Cluster", cluster).Msg("Failed Hosts:")
		PrintCommand(fmt.Sprintf("cat %v", getClusterFile("fail")))
	})
}

func (p *Pssh) RunRange(cmd string, cluster string, parallelism int, disableOutput bool, start int, end int) {
	if start != -1 && end != -1 {
		subClusterName := cluster + "m"
		ExtractSubCluster(cluster, subClusterName, start-1, end)
		p.Run(cmd, subClusterName, parallelism, disableOutput)
	} else {
		p.Run(cmd, cluster, parallelism, disableOutput)
	}
}

func clearOutputPaths() {
	util.ClearDirectory(config.OUTPUT_PATH)
	util.ClearDirectory(config.ERROR_PATH)
}

func ExtractSubCluster(clusterName string, subClusterName string, start int, end int) {
	ips := util.ReadAllLines(getClusterFile(clusterName))
	WriteClusterFile(subClusterName, strings.Join(ips[start:end], "\n"))
}

func WriteClusterFile(clusterName string, content string) {
	filePath := getClusterFile(clusterName)
	if err := os.WriteFile(filePath, []byte(content), util.DEFAULT_PERM); err != nil {
		log.Error().Err(err).Str("Cluster", clusterName).Str("Path", filePath).Msg("Failed to write cluster file")
	}
}

func ReadClusterFile(clusterName string) []string {
	filePath := getClusterFile(clusterName)
	return util.ReadAllLines(filePath)
}

func RemoveCluster(mainClusterName string, removeClusterName string) int {
	mainSet := ReadClusterFile(mainClusterName)
	removeSet := ReadClusterFile(removeClusterName)
	finalSet := lo.Without(mainSet, removeSet...)
	WriteClusterFile(mainClusterName, strings.Join(finalSet, "\n"))
	return len(mainSet) - len(finalSet)
}

func GetClusterHost(clusterName string, index int) string {
	ips := ReadClusterFile(clusterName)
	if index > len(ips) {
		return "INVALID"
	}
	return ips[index-1]
}

func SearchCluster(keyword string) (clusters []string) {
	log.Info().Str("Path", config.CLUSTER_PATH).Msg("Searching")
	files, _ := filepath.Glob(fmt.Sprintf("%v/*%v*", config.CLUSTER_PATH, keyword))
	for _, name := range files {
		fileName := strings.Replace(name, config.CLUSTER_PATH+"/", "", 1)
		cluster := strings.TrimSuffix(fileName, ".txt")
		clusters = append(clusters, cluster)
	}
	return
}

func Md5Checker(cmd string, cluster string) {
	/* Run Command to get Ip Wise output */
	files := runMd5Command(cmd, cluster)
	hashMap, sortList := computeMd5Hashes(files)
	analyzeMd5Results(cmd, cluster, hashMap, sortList)
}

func runMd5Command(cmd string, cluster string) map[string][]string {
	FastPssh.Run(cmd, cluster, defaultTimeout, true)
	return util.ReadFileMap(config.OUTPUT_PATH, true)
}

/* Compute Md5 and store as list with count */
func computeMd5Hashes(files map[string][]string) (map[string]*util.Md5Info, []*util.Md5Info) {
	hashMap := map[string]*util.Md5Info{}
	var sortList []*util.Md5Info

	for path, content := range files {
		md5Hash := util.GetMD5Hash(strings.Join(content, "\n"))
		if _, ok := hashMap[md5Hash]; !ok {
			info := &util.Md5Info{FileList: []string{}, Hash: md5Hash}
			hashMap[md5Hash] = info
			sortList = append(sortList, info)
		}
		hashMap[md5Hash].Add(path)
	}
	return hashMap, sortList
}

func analyzeMd5Results(cmd string, cluster string, hashMap map[string]*util.Md5Info, sortList []*util.Md5Info) {
	if len(sortList) > 1 {
		logMultipleMd5(cmd, cluster, sortList)
		compareMd5Results(cluster, sortList)
	} else {
		log.Info().Str("Cluster", cluster).Str("Hash", sortList[0].Hash).Str("Count", fmt.Sprint(sortList[0].Count)).Msg("Cluster Homogenous")
	}
}

func logMultipleMd5(cmd string, cluster string, sortList []*util.Md5Info) {
	log.Warn().Str("Cluster", cluster).Str("CMD", cmd).Msg("Multiple MD5 Detected")

	/* Sort Md5 List by Count */
	sort.Slice(sortList, func(i, j int) bool {
		return sortList[i].Count > sortList[j].Count
	})
	for _, value := range sortList {
		log.Warn().Str("Cluster", cluster).Str("Hash", value.Hash).Int("Count", value.Count).Msg("MD5 Result")
	}
}

func compareMd5Results(cluster string, sortList []*util.Md5Info) {
	first := sortList[0]
	firstFile := first.FileList[0]
	for i := 1; i < len(sortList); i++ {
		current := sortList[i]
		currentFile := current.FileList[0]
		log.Info().Str("First", firstFile).Str("FirstHash", first.Hash).Str("Current", currentFile).Str("CurrentHash", current.Hash).Msg("Diffing")
		if util.IsDebugMode() {
			util.PrintFile(firstFile, firstFile)
			util.PrintFile(currentFile, currentFile)
		}
		fmt.Println(RunCommandIgnoreError(fmt.Sprintf("colordiff %v %v", firstFile, currentFile)))
	}
}

func SearchContent(regex string) string {
	return RunCommandIgnoreError(fmt.Sprintf("grep -inrR '%v' %v", regex, config.OUTPUT_PATH))
}

func getClusterFile(name string) string {
	return fmt.Sprintf("%v/%v.txt", config.CLUSTER_PATH, name)
}

func (p *Pssh) getDisplayFlag() string {
	if !p.displayOutput {
		return ""
	}
	return "-P"
}
