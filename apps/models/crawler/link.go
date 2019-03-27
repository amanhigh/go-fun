package crawler

import (
	"fmt"

	"errors"

	"github.com/amanhigh/go-fun/util"
)

type LinkInfo struct {
	Link string
}

func (self *LinkInfo) Print() {
	util.PrintWhite(fmt.Sprintf("%v", self.Link))
}

func (self *LinkInfo) GoodBad() (err error) {
	if self.Link == "" {
		err = errors.New("Bad Link: " + self.Link)
	}
	return
}

func (self *LinkInfo) ToUrl() []string {
	return []string{self.Link}
}
