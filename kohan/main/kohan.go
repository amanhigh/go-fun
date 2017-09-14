package main

import (
	"fmt"
	"os"
	"github.com/amanhigh/go-fun/kohan/processor"
	"github.com/amanhigh/go-fun/command"
	processor2 "github.com/Flipkart/elb/scripts/kohan/processor"
)

var PROCESSOR_MAP = map[string]processor.ProcessorI{
	"expose": &processor.Processor{&processor.ExposeProcessor{}},
	"elb": &processor.Processor{&processor2.ElbProcessor{}},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: kohan <Processor Name> <Command Name>")
		os.Exit(1)
	} else if len(os.Args) < 3 {
		processorName := os.Args[1]
		if p, ok := PROCESSOR_MAP[processorName]; ok {
			commander.PrintWhite(p.Help())
		} else {
			commander.PrintRed("Unknown Processor: " + processorName)
			commander.PrintWhite("Valid List")
		}
		os.Exit(1)
	}

	processorName := os.Args[1]
	command := os.Args[2]

	if p, ok := PROCESSOR_MAP[processorName]; ok {
		p.Process(command, os.Args[3:])
	} else {
		commander.PrintRed("Unknown Processor: " + processorName)
		commander.PrintWhite("Valid List")
	}

}

func Help() {
	//	"cosmosd": &processor2.CosmosDebugProcessor{p},
}
