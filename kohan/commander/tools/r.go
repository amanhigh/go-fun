package tools

import "fmt"

var R_STAT_FILE = "/tmp/r.stat"

func StatSummary(columnIndex int, separator string) {
	PrintCommand(fmt.Sprintf(`R --slave -e 'x <- read.table(file="%v",sep="%v"); summary(x[%v]); quantile(x[,%v],c(0.1,0.25,0.5,0.9,0.95,0.99,0.999))'`, R_STAT_FILE, separator, columnIndex, columnIndex))
}
