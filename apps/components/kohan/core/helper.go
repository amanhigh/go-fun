package core

import (
	"fmt"
	util2 "github.com/amanhigh/go-fun/apps/common/util"
	"strconv"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/util"
)

func Printf(templateFile string, paramFile string, marker string) {

	for _, templateSplit := range util.ReadAllLines(templateFile) {
		for _, paramLine := range util.ReadAllLines(paramFile) {

			template := templateSplit
			for i, param := range strings.Split(paramLine, " ") {
				template = strings.Replace(template, marker+strconv.Itoa(i+1), param, -1)
			}
			fmt.Println(template)
		}
		fmt.Println("")
	}
}

/**
Reformat Lines as per ForexTester 3 Format.
Date, Open, High, Low, Close, Volume
*/
func ReformatInvestingFile(filePath string) (err error) {
	//Read all Lines
	lines := util.ReadAllLines(filePath)
	var outLines []string
	for i, line := range lines {
		//Skip Header
		if i > 0 {
			//Split Fields
			split := strings.Split(line, `","`)

			//Remove Internal Commas
			for i, word := range split {
				split[i] = strings.Replace(word, ",", "", -1)
			}

			//Parse Date
			var formattedLine []string
			if parse, err := time.Parse("\"Jan 02 2006", split[0]); err == nil {
				formattedDate := parse.Format("20060102")
				//Add Date First
				formattedLine = append(formattedLine, formattedDate)
				//Add Open,High,Low
				formattedLine = append(formattedLine, split[2:5]...)
				//Add Close, Volume
				formattedLine = append(formattedLine, split[1], split[5])

				join := strings.Join(formattedLine, ",")
				//fmt.Println(join)
				outLines = append(outLines, join)
			} else {
				return err
			}
		}
	}

	//Data from oldest to Newest
	util2.ReverseArray(outLines)
	//OverWrite Output Lines
	return util.WriteLines(filePath, outLines)
}
