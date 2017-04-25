package main

import (
	"fmt"
	"os"
	"os/exec"
	"io/ioutil"
	"strings"
	"strconv"
)

func main() {
	args := os.Args[1:]
	if (len(args) < 2) {
		fmt.Println("Usage: sship <Cluster> <No.>")
	} else {
		cluster := args[0]
		n, _ := strconv.Atoi(args[1])

		filePath := fmt.Sprintf("/tmp/clusters/%v.txt", cluster)

		ips := readFile(filePath)

		l := len(ips)
		if (l < n) {
			fmt.Printf("Insufficient Range: %v < %v - %v\n", l, n, cluster)
		} else {
			ip := ips[n]
			fmt.Printf("Entering IP: %s", ip)
			ssh(ip)
		}
	}
}

func ssh(ip string) {
	execCmd(fmt.Sprintf("ssh %s", ip))
}

func execCmd(cmd string) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err == nil {
		fmt.Printf("%s", out)
	} else {
		fmt.Printf("Error Occured: %s", err)
	}
}
func readFile(filePath string) []string {
	content, err := ioutil.ReadFile(filePath)
	if (err == nil) {
		return strings.Split(string(content), "\n")
	}
	return nil
}