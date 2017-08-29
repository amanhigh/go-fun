package main

import (
	"fmt"
	"os"
	"github.com/amanhigh/go-fun/kohan/processor"
	processor2 "github.com/Flipkart/elb/scripts/kohan/processor"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: kohan <Processor Name> <Command Name>")
		os.Exit(1)
	}

	processorName := os.Args[1]
	command := os.Args[2]

	pMap := getProcessorMap()
	if p, ok := pMap[processorName]; ok {
		p.Process(command)
	} else {
		fmt.Println("Unknown Processor:", processorName)
	}

}

func getProcessorMap() map[string]processor.ProcessorI {
	p := processor.Processor{Args: os.Args[3:]}

	return map[string]processor.ProcessorI{
		"ssh":    &processor.SshProcessor{p},
		"expose": &processor.ExposeProcessor{p},
		"elb":    &processor2.ElbProcessor{p},
	}
}
