package crawler

import (
	"fmt"

	"github.com/amanhigh/go-fun/util"
	"github.com/amanhigh/go-fun/util/helper"
)

type GameInfo struct {
	Name string
	Link string
}

func (self *GameInfo) Print() {
	util.PrintWhite(fmt.Sprintf("%v: %v", self.Name, self.Link))
}

func (self *GameInfo) GoodBad() error {
	return nil
}

func (self *GameInfo) ToUrl() []string {
	return []string{self.Link, helper.YoutubeSearch(self.Name + " Review")}
}
