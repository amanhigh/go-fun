package processor

import (
	"flag"
	"fmt"
	"strings"
	"github.com/amanhigh/go-fun/commander"
)

/* Interface */
type HandleFunc func(*flag.FlagSet, []string) error
type DirectFunc func()

type HandlerI interface {
	GetArgedHandlers() (map[string]HandleFunc)
	GetNonArgedHandlers() (map[string]DirectFunc)
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

	if handleFunc, ok := self.Handler.GetArgedHandlers()[commandName]; ok {
		handleFunc(flagSet, args)
	} else if directFunc, ok := self.Handler.GetNonArgedHandlers()[commandName]; ok {
		directFunc()
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
	return fmt.Sprintf("Commands:\n Flagged - %v\n Direct - %v\n", getCommandHelpString(self.Handler.GetArgedHandlers()), getDirectHelpString(self.Handler.GetNonArgedHandlers()))
}

func getCommandHelpString(funcs map[string]HandleFunc) (string) {
	commandNames := []string{}
	for command := range funcs {
		commandNames = append(commandNames, command)
	}
	return strings.Join(commandNames, ", ")
}
func getDirectHelpString(funcs map[string]DirectFunc) (string) {
	commandNames := []string{}
	for command := range funcs {
		commandNames = append(commandNames, command)
	}
	return strings.Join(commandNames, ", ")
}
