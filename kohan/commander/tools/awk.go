package tools

import (
	"time"
	"fmt"
	"github.com/amanhigh/go-fun/util"
)

var TIME_OUT_FILE = "/tmp/time.out"

func ExtractLogForDuration(remoteIp string, logFile string, grepPattern string, startTime time.Time, endTime time.Time, timeLayout string) {
	startString := util.FormatTime(startTime, timeLayout)
	endString := util.FormatTime(endTime, timeLayout)
	util.PrintYellow(fmt.Sprintf("Extracting Log between %v - %v File:%v Server:%v Output:%v", startString, endString, logFile, remoteIp, TIME_OUT_FILE))
	PrintCommand(fmt.Sprintf(`ssh %v "awk '\$2>=\"%v\" && \$2<=\"%v\"' %v | grep '%v' > %v"`,
		remoteIp, startString, endString, logFile, grepPattern, TIME_OUT_FILE))
}
