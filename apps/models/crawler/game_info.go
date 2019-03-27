package crawler

import (
	"fmt"

	"errors"

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

func (self *GameInfo) GoodBad() (err error) {
	if self.Name == "" || self.Link == "" {
		err = errors.New("Bad Game: " + self.Name)
	}
	return
}

func (self *GameInfo) ToUrl() []string {
	return []string{self.Link, helper.YoutubeSearch(self.Name + " Review")}
}
