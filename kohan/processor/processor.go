package processor

import (
	"flag"
	"fmt"
	"github.com/amanhigh/go-fun/command"
	"strings"
)

/* Interface */
type HandleFunc func(*flag.FlagSet, []string) error

type HandlerI interface {
	GetHandleMap() (map[string]HandleFunc)
}

type ProcessorI interface {
	Process(string, []string) (bool)
	Help() string
}

/* Abstract Processor */
type Processor struct {
	Handler HandlerI
}

func (self *Processor) Process(commandName string, args []string) (bool) {
	var e error
	flagSet := flag.NewFlagSet(commandName, flag.ExitOnError)

	handlerMap := self.Handler.GetHandleMap()
	if handleFunc, ok := handlerMap[commandName]; ok {
		handleFunc(flagSet, args)
	} else {
		commander.PrintWhite(self.Help())
	}

	if e != nil {
		fmt.Println(e.Error())
		flagSet.Usage()
		return false
	}
	return true
}

func (self *Processor) Help() string {
	commandNames := []string{}
	for command := range self.Handler.GetHandleMap() {
		commandNames = append(commandNames, command)
	}
	return "Commands: " + strings.Join(commandNames, ", ")
}
