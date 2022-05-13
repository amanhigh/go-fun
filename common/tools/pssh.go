package tools

import (
	"fmt"
	util2 "github.com/amanhigh/go-fun/common/util"
	config2 "github.com/amanhigh/go-fun/models/config"
	"github.com/fatih/color"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

var FastPssh = Pssh{20, config2.OUTPUT_PATH, config2.ERROR_PATH, false}
var NORMAL_PSSH = Pssh{30, config2.OUTPUT_PATH, config2.ERROR_PATH, false}
var DisplayPssh = Pssh{10, config2.OUTPUT_PATH, config2.ERROR_PATH, true}
var SlowPssh = Pssh{240, config2.OUTPUT_PATH, config2.ERROR_PATH, false}

type Pssh struct {
	Timeout       int
	outputPath    string
	errorPath     string
	displayOutput bool
}

func (self *Pssh) Run(cmd string, cluster string, parallelism int, disableOutput bool) {
	clearOutputPaths()

	psshCmd := fmt.Sprintf(`script %v pssh -h %v -t %v -o %v -e %v %v -p %v '%v'`,
		config2.CONSOLE_FILE, getClusterFile(cluster), self.Timeout, self.outputPath, self.errorPath, self.getDisplayFlag(), parallelism, cmd)
	if disableOutput {
		RunCommandPrintError(psshCmd)
	} else {
		color.White(fmt.Sprintf("Running Parallel SSH. Cluster: %v Parallelism:%v", cluster, parallelism))
		LiveCommand(psshCmd)
	}

	RunIf(fmt.Sprintf("grep FAILURE %v", getClusterFile("console")), func(output string) {
		PrintCommand(fmt.Sprintf("grep SUCCESS %v | awk '{print $4}' > %v", getClusterFile("console"), getClusterFile("pass")))
		PrintCommand(fmt.Sprintf("grep FAILURE %v | awk '{print $4}' > %v", getClusterFile("console"), getClusterFile("fail")))
		color.Yellow("Failed Hosts:")
		PrintCommand(fmt.Sprintf("cat %v", getClusterFile("fail")))
	})
}

func (self *Pssh) RunRange(cmd string, cluster string, parallelism int, disableOutput bool, start int, end int) {
	if start != -1 && end != -1 {
		subClusterName := cluster + "m"
		ExtractSubCluster(cluster, subClusterName, start-1, end)
		self.Run(cmd, subClusterName, parallelism, disableOutput)
	} else {
		self.Run(cmd, cluster, parallelism, disableOutput)
	}
}

func clearOutputPaths() {
	util2.ClearDirectory(config2.OUTPUT_PATH)
	util2.ClearDirectory(config2.ERROR_PATH)
}

func ExtractSubCluster(clusterName string, subClusterName string, start int, end int) {
	ips := util2.ReadAllLines(getClusterFile(clusterName))
	WriteClusterFile(subClusterName, strings.Join(ips[start:end], "\n"))
}

func WriteClusterFile(clusterName string, content string) {
	filePath := getClusterFile(clusterName)
	ioutil.WriteFile(filePath, []byte(content), util2.DEFAULT_PERM)
}

func ReadClusterFile(clusterName string) []string {
	filePath := getClusterFile(clusterName)
	return util2.ReadAllLines(filePath)
}

func RemoveCluster(mainClusterName string, removeClusterName string) int {
	mainSet := ReadClusterFile(mainClusterName)
	removeSet := ReadClusterFile(removeClusterName)
	diff, _ := funk.Difference(mainSet, removeSet)
	finalSet := diff.([]string)
	WriteClusterFile(mainClusterName, strings.Join(finalSet, "\n"))
	return len(mainSet) - len(finalSet)
}

func GetClusterHost(clusterName string, index int) string {
	ips := ReadClusterFile(clusterName)
	if index <= len(ips) {
		return ips[index-1]
	} else {
		return "INVALID"
	}
}

func SearchCluster(keyword string) (clusters []string) {
	color.Blue("Searching: " + config2.CLUSTER_PATH)
	files, _ := filepath.Glob(fmt.Sprintf("%v/*%v*", config2.CLUSTER_PATH, keyword))
	for _, name := range files {
		fileName := strings.Replace(name, config2.CLUSTER_PATH+"/", "", 1)
		cluster := strings.TrimRight(fileName, ".txt")
		clusters = append(clusters, cluster)
	}
	return
}

func Md5Checker(cmd string, cluster string) {
	/* Run Command to get Ip Wise output */
	FastPssh.Run(cmd, cluster, 200, true)
	files := util2.ReadFileMap(config2.OUTPUT_PATH, true)

	/* Compute Md5 and store as list with count */
	hashMap := map[string]*util2.Md5Info{}
	var sortList []*util2.Md5Info

	for path, content := range files {
		md5Hash := util2.GetMD5Hash(strings.Join(content, "\n"))
		if _, ok := hashMap[md5Hash]; !ok {
			info := &util2.Md5Info{FileList: []string{}, Hash: md5Hash}
			hashMap[md5Hash] = info
			sortList = append(sortList, info)
		}
		hashMap[md5Hash].Add(path)
	}

	/* If more than one Md5 Sums Found */
	if len(sortList) > 1 {
		color.Red("Multiple MD5 Detected, Cluster Non Homogenous: %v", cluster)

		/* Sort Md5 List by Count */
		sort.Slice(sortList, func(i, j int) bool {
			return sortList[i].Count > sortList[j].Count
		})
		for _, value := range sortList {
			color.Blue("%v %v", value.Hash, value.Count)
		}

		/* Perform Diff on first file of top two md5's */
		first := sortList[0]
		firstFile := first.FileList[0]
		for i := 1; i < len(sortList); i++ {
			current := sortList[i]
			currentFile := current.FileList[0]
			color.Cyan("Diffing Top with Current: %v (%v) vs %v (%v)", firstFile, first.Hash, currentFile, current.Hash)
			if util2.IsDebugMode() {
				util2.PrintFile(firstFile, firstFile)
				util2.PrintFile(currentFile, currentFile)
			}
			fmt.Println(RunCommandIgnoreError(fmt.Sprintf("colordiff %v %v", firstFile, currentFile)))
		}
	} else {
		color.Green(fmt.Sprintf("Single Md5 Found, Cluster Homogenous: %v Hash:%v Count:%v", cluster, sortList[0].Hash, sortList[0].Count))
	}
}

func SearchContent(regex string) string {
	return RunCommandIgnoreError(fmt.Sprintf("grep -inrR '%v' %v", regex, config2.OUTPUT_PATH))
}

func getClusterFile(name string) string {
	return fmt.Sprintf("%v/%v.txt", config2.CLUSTER_PATH, name)
}

func (self *Pssh) getDisplayFlag() string {
	if self.displayOutput {
		return "-P"
	} else {
		return ""
	}
}
