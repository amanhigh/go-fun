package crawler

import (
	"fmt"
	"github.com/fatih/color"

	"errors"
)

type LinkInfo struct {
	Link string
}

func (self *LinkInfo) Print() {
	color.White(fmt.Sprintf("%v", self.Link))
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
