package crawler

import (
	"github.com/amanhigh/go-fun/util"
	"fmt"
)

type ImdbInfo struct {
	Name   string
	Link   string
	Rating float64
	CutOff int
}

func (self *ImdbInfo) Print() {
	util.PrintWhite(fmt.Sprintf("%v: %.2f - %v", self.Name, self.Rating, self.Link))
}

func (self *ImdbInfo) GoodBad() bool {
	return self.Rating >= float64(self.CutOff) || self.Rating < 0.1
}

func (self *ImdbInfo) ToUrl() string {
	return self.Link
}
