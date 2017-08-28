package commander

import "fmt"

var FastPssh = Pssh{
	10,
	OUTPUT_PATH,
	ERROR_PATH,
	false,
}

type Pssh struct {
	Timeout       int
	outputPath    string
	errorPath     string
	displayOutput bool
}

func (self *Pssh) Run(cmd string, cluster string, parallelism int) {
	psshCmd := fmt.Sprintf("script %v pssh -h %v -t %v -o %v -e %v %v -p %v %v;",
		CONSOLE_FILE, getClusterFile(cluster), self.Timeout, self.outputPath, self.errorPath, self.getDisplayFlag(), parallelism, cmd)
	PrintCommand(psshCmd)
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
