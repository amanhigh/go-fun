package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/amanhigh/go-fun/apps/components/kohan/core"
	"github.com/amanhigh/go-fun/apps/models/config"
	"github.com/fatih/color"
	"sort"
	"strings"

	"github.com/amanhigh/go-fun/apps/common/tools"
	"github.com/amanhigh/go-fun/util"
)

type md5Info struct {
	hash     string
	fileList []string
	count    int
}

func (self *md5Info) Add(path string) {
	self.fileList = append(self.fileList, path)
	self.count++
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Md5Checker(cmd string, cluster string) {
	/* Run Command to get Ip Wise output */
	tools.FastPssh.Run(cmd, cluster, 200, true)
	files := util.ReadFileMap(config.OUTPUT_PATH, true)

	/* Compute Md5 and store as list with count */
	hashMap := map[string]*md5Info{}
	sortList := []*md5Info{}

	for path, content := range files {
		md5Hash := GetMD5Hash(strings.Join(content, "\n"))
		if _, ok := hashMap[md5Hash]; !ok {
			info := &md5Info{fileList: []string{}, hash: md5Hash}
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
			return sortList[i].count > sortList[j].count
		})
		for _, value := range sortList {
			color.Blue("%v %v", value.hash, value.count)
		}

		/* Perform Diff on first file of top two md5's */
		first := sortList[0]
		firstFile := first.fileList[0]
		for i := 1; i < len(sortList); i++ {
			current := sortList[i]
			currentFile := current.fileList[0]
			color.Cyan("Diffing Top with Current: %v (%v) vs %v (%v)", firstFile, first.hash, currentFile, current.hash)
			if core.IsDebugMode() {
				util.PrintFile(firstFile, firstFile)
				util.PrintFile(currentFile, currentFile)
			}
			fmt.Println(tools.RunCommandIgnoreError(fmt.Sprintf("colordiff %v %v", firstFile, currentFile)))
		}
	} else {
		color.Green(fmt.Sprintf("Single Md5 Found, Cluster Homogenous: %v Hash:%v Count:%v", cluster, sortList[0].hash, sortList[0].count))
	}
}
