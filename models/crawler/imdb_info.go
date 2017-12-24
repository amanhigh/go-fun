package crawler

import (
	"github.com/amanhigh/go-fun/util"
	"fmt"
	"errors"
)

type ImdbInfo struct {
	Name     string
	Link     string
	Rating   float64
	MyRating float64
	CutOff   int
}

func (self *ImdbInfo) Print() {
	util.PrintWhite(fmt.Sprintf("%v: %.2f/%.2f - %v", self.Name, self.MyRating, self.Rating, self.Link))
}

func (self *ImdbInfo) GoodBad() error {
	if self.Rating < float64(self.CutOff) {
		return errors.New(fmt.Sprintf("Subpar Rating %v < %v", self.Rating, self.CutOff))
	} else if self.MyRating != -1 {
		return errors.New(fmt.Sprintf("Movie already Rated: %v", self.MyRating))
	}
	return nil
}

func (self *ImdbInfo) ToUrl() string {
	return self.Link
}
