package crawler

import (
	"github.com/amanhigh/go-fun/util"
	"fmt"
)

type ImdbInfo struct {
	Name   string
	Link   string
	Rating float64
}

func (self *ImdbInfo) Print() {
	util.PrintWhite(fmt.Sprintf("%v: %.2f - %v", self.Name, self.Rating, self.Link))
}
