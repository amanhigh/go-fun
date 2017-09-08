package commander

import (
	"fmt"
	"io/ioutil"
)

var FastPssh = Pssh{20, OUTPUT_PATH, ERROR_PATH, false,}
var NORMAL_PSSH=Pssh{30, OUTPUT_PATH, ERROR_PATH, false,}
var DisplayPssh = Pssh{10, OUTPUT_PATH, ERROR_PATH, true,}
var SlowPssh = Pssh{240, OUTPUT_PATH, ERROR_PATH, false,}

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

	RunIf(fmt.Sprintf("grep FAILURE %v", getClusterFile("console.txt")), func(output string) {
		PrintCommand(fmt.Sprintf("grep FAILURE %v | awk '{print $4}' > %v", getClusterFile("console"), getClusterFile("fail")))
		PrintYellow("Failed Hosts:")
		PrintCommand(fmt.Sprintf("cat %v", getClusterFile("fail")))
	})
}
func clearOutputPaths() {
	ClearDirectory(OUTPUT_PATH)
	ClearDirectory(ERROR_PATH)
}

func WriteClusterFile(clusterName string, content string) {
	filePath := getClusterFile(clusterName)
	ioutil.WriteFile(filePath, []byte(content), DEFAULT_PERM)
}

func ReadClusterFile(clusterName string) []string {
	filePath := getClusterFile(clusterName)
	return ReadAllLines(filePath)
}

func IndexedIp(clusterName string, index int) {
	ips := ReadClusterFile(clusterName)
	if index <= len(ips) {
		fmt.Println(ips[index-1])
	} else {
		fmt.Println("INVALID")
	}
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
