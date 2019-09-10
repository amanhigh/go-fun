package tools

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	. "github.com/amanhigh/go-fun/util"
)

var FastPssh = Pssh{20, OUTPUT_PATH, ERROR_PATH, false}
var NORMAL_PSSH = Pssh{30, OUTPUT_PATH, ERROR_PATH, false}
var DisplayPssh = Pssh{10, OUTPUT_PATH, ERROR_PATH, true}
var SlowPssh = Pssh{240, OUTPUT_PATH, ERROR_PATH, false}

type Pssh struct {
	Timeout       int
	outputPath    string
	errorPath     string
	displayOutput bool
}

func (self *Pssh) Run(cmd string, cluster string, parallelism int, disableOutput bool) {
	clearOutputPaths()

	psshCmd := fmt.Sprintf(`script %v pssh -h %v -t %v -o %v -e %v %v -p %v '%v'`,
		CONSOLE_FILE, getClusterFile(cluster), self.Timeout, self.outputPath, self.errorPath, self.getDisplayFlag(), parallelism, cmd)
	if disableOutput {
		RunCommandPrintError(psshCmd)
	} else {
		PrintWhite(fmt.Sprintf("Running Parallel SSH. Cluster: %v Parallelism:%v", cluster, parallelism))
		LiveCommand(psshCmd)
	}

	RunIf(fmt.Sprintf("grep FAILURE %v", getClusterFile("console")), func(output string) {
		PrintCommand(fmt.Sprintf("grep SUCCESS %v | awk '{print $4}' > %v", getClusterFile("console"), getClusterFile("pass")))
		PrintCommand(fmt.Sprintf("grep FAILURE %v | awk '{print $4}' > %v", getClusterFile("console"), getClusterFile("fail")))
		PrintYellow("Failed Hosts:")
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
	ClearDirectory(OUTPUT_PATH)
	ClearDirectory(ERROR_PATH)
}

func ExtractSubCluster(clusterName string, subClusterName string, start int, end int) {
	ips := ReadAllLines(getClusterFile(clusterName))
	WriteClusterFile(subClusterName, strings.Join(ips[start:end], "\n"))
}

func WriteClusterFile(clusterName string, content string) {
	filePath := getClusterFile(clusterName)
	ioutil.WriteFile(filePath, []byte(content), DEFAULT_PERM)
}

func ReadClusterFile(clusterName string) []string {
	filePath := getClusterFile(clusterName)
	return ReadAllLines(filePath)
}

func RemoveCluster(mainClusterName string, removeClusterName string) int {
	mainSet := ReadClusterFile(mainClusterName)
	removeSet := ReadClusterFile(removeClusterName)
	finalSet := SliceMinus(mainSet, removeSet)
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
	PrintBlue("Searching: " + CLUSTER_PATH)
	files, _ := filepath.Glob(fmt.Sprintf("%v/*%v*", CLUSTER_PATH, keyword))
	for _, name := range files {
		fileName := strings.TrimLeft(name, CLUSTER_PATH)
		cluster := strings.TrimRight(fileName, ".txt")
		clusters = append(clusters, cluster)
	}
	return
}

func SearchContent(regex string) string {
	return RunCommandIgnoreError(fmt.Sprintf("grep -inrR '%v' %v", regex, OUTPUT_PATH))
}

func getClusterFile(name string) string {
	return fmt.Sprintf("%v/%v.txt", CLUSTER_PATH, name)
}

func (self *Pssh) getDisplayFlag() string {
	if self.displayOutput {
		return "-P"
	} else {
		return ""
	}
}
