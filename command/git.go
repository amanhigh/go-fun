package commander

import "fmt"

func GetBranch() string {
	output, _ := runCommand("git rev-parse --abbrev-ref HEAD")
	return output
}

func GitCommit(msg string, filePath string) {
	PrintCommand(fmt.Sprintf("git commit -m '%v' %v", msg, filePath))
}

func GitPush(){
	PrintCommand("git push origin")
}

func GitDiff() {
	PrintCommand("git diff -U0")
}
