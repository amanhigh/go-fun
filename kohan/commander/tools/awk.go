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
	util.PrintYellow(fmt.Sprintf("Extracting Log between %v - %v FilterPattern:%v File:%v Server:%v Output:%v", startString, endString, grepPattern, logFile, remoteIp, TIME_OUT_FILE))
	PrintCommand(fmt.Sprintf(`ssh %v "cat %v | awk '\$2>=\"%v\" && \$2<=\"%v\"' | grep '%v' > %v"`,
		remoteIp, logFile, startString, endString, grepPattern, TIME_OUT_FILE))
}

func PipeForDuration(remoteIp string, logFile string, grepPattern string, startTime time.Time, endTime time.Time, timeLayout string, pipe string) {
	startString := util.FormatTime(startTime, timeLayout)
	endString := util.FormatTime(endTime, timeLayout)
	util.PrintYellow(fmt.Sprintf("Filtering between %v - %v FilterPattern:%v File:%v Server:%v Pipe:%v", startString, endString, grepPattern, logFile, remoteIp, pipe))
	PrintCommand(fmt.Sprintf(`ssh %v "cat %v | awk '\$2>=\"%v\" && \$2<=\"%v\"' | grep '%v' | %v"`, remoteIp, logFile, startString, endString, grepPattern, pipe))
}
