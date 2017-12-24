package crawler

import (
	"github.com/amanhigh/go-fun/util"
	"fmt"
	"errors"
	"github.com/amanhigh/go-fun/util/helper"
)

type ImdbInfo struct {
	Name     string
	Link     string
	Language string
	Rating   float64
	MyRating float64
	CutOff   int
}

func (self *ImdbInfo) Print() {
	if self.MyRating != -1 {
		util.PrintWhite(fmt.Sprintf("%v: %.2f/%.2f - %v", self.Name, self.MyRating, self.Rating, self.Link))
	} else {
		util.PrintWhite(fmt.Sprintf("%v: %.2f - %v", self.Name, self.Rating, self.Link))
	}
}

func (self *ImdbInfo) GoodBad() error {
	if self.Rating < float64(self.CutOff) {
		return errors.New(fmt.Sprintf("Subpar Rating %v < %v", self.Rating, self.CutOff))
	} else if self.MyRating != -1 {
		return errors.New(fmt.Sprintf("Movie already Rated: %v", self.MyRating))
	}
	return nil
}

func (self *ImdbInfo) ToUrl() []string {
	return []string{self.Link, helper.YoutubeSearch(self.Name), self.getDownloadLink()}
}

func (self *ImdbInfo) getDownloadLink() string {
	switch self.Language {
	case "en":
		return helper.YtsSearch(self.Name)
	case "hi":
		return helper.HotStarSearch(self.Name)
	default:
		return helper.YoutubeSearch(self.Name + " Full Movie")
	}

}
