package processor

import (
	"flag"
	"fmt"
	"github.com/amanhigh/go-fun/command"
)

type ExposeProcessor struct {
	Processor
}

func (p *ExposeProcessor) Process(commandName string) (bool) {
	var e error
	flagSet := flag.NewFlagSet(commandName, flag.ExitOnError)

	switch commandName {
	case "pssh":
		cmd := flagSet.String("cmd", "", "Command To Run")
		cluster := flagSet.String("cl", "", "Cluster To Run On")
		parallelism := flagSet.Int("p", 50, "Parallelism")
		e = flagSet.Parse(p.Args)
		commander.FastPssh.Run(*cmd, *cluster, *parallelism)
	default:
		fmt.Println(p.Help())
	}

	if e != nil {
		fmt.Println(e.Error())
		flagSet.Usage()
		return false
	}
	return true
}

func (p *ExposeProcessor) Help() string {
	return `Commands: pssh`
}
