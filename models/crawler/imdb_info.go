package crawler

import (
	"github.com/amanhigh/go-fun/util"
	"fmt"
	"errors"
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

func (self *ImdbInfo) GoodBad() error {
	if self.Rating >= float64(self.CutOff) || self.Rating < 0.1 {
		return errors.New(fmt.Sprintf("Subpar Rating %v < %v", self.Rating, self.CutOff))
	}
	return nil
}

func (self *ImdbInfo) ToUrl() string {
	return self.Link
}
