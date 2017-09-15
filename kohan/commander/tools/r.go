package tools

import "fmt"

var R_STAT_FILE = "/tmp/r.stat"

func StatSummary(columnIndex int, separator string) {
	PrintCommand(fmt.Sprintf(`R --slave -e 'x <- read.table(file="%v",sep="%v"); summary(x[%v]); quantile(x[,%v])'`, R_STAT_FILE, separator, columnIndex, columnIndex))
}
