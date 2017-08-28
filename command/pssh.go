package commander

import "fmt"

var FastPssh = Pssh{10, OUTPUT_PATH, ERROR_PATH, false,}
var DisplayPssh = Pssh{10, OUTPUT_PATH, ERROR_PATH, true,}
var SlowPssh = Pssh{240, OUTPUT_PATH, ERROR_PATH, false,}

type Pssh struct {
	Timeout       int
	outputPath    string
	errorPath     string
	displayOutput bool
}

func (self *Pssh) Run(cmd string, cluster string, parallelism int) {
	psshCmd := fmt.Sprintf("script %v pssh -h %v -t %v -o %v -e %v %v -p %v '%v';",
		CONSOLE_FILE, getClusterFile(cluster), self.Timeout, self.outputPath, self.errorPath, self.getDisplayFlag(), parallelism, cmd)
	PrintCommand(psshCmd)

	RunIf(fmt.Sprintf("grep FAILURE %v", getClusterFile("console.txt")), func() {
		PrintCommand(fmt.Sprintf("grep FAILURE %v | awk '{print $4}' > %v", getClusterFile("console"), getClusterFile("fail")))
		PrintYellow("Failed Hosts:")
		PrintCommand(fmt.Sprintf("cat %v", getClusterFile("fail")))
	})

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
