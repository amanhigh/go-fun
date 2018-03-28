package main

import (
	"fmt"
	"os"
	"github.com/amanhigh/go-fun/kohan/processor"
	processor2 "github.fkinternal.com/Flipkart/elb/deployment/scripts/kohan/processor"
	"strings"
	"github.com/amanhigh/go-fun/util"
)

var PROCESSOR_MAP = map[string]processor.ProcessorI{
	"expose":  &processor.Processor{&processor.ExposeProcessor{}},
	"crawl":  &processor.Processor{&processor.CrawlProcessor{}},
	"elb":     &processor.Processor{&processor2.ElbProcessor{}},
	"infra":     &processor.Processor{&processor2.InfraProcessor{}},
	"cosmosd": &processor.Processor{&processor2.CosmosDebugProcessor{}},
}

func main() {
	/* Processor Not Specified */
	if len(os.Args) < 2 {
		fmt.Println("Usage: kohan <Processor Name> <Command Name>")
		Help()
		os.Exit(1)
	}

	processorName := os.Args[1]
	selectedProcessor, ok := PROCESSOR_MAP[processorName]
	/* Specified Processor Not Found */
	if !ok {
		util.PrintRed("Unknown Processor: " + processorName)
		Help()
		os.Exit(1)
	}

	/* Command Not Specified */
	if len(os.Args) < 3 {
		util.PrintWhite(selectedProcessor.Help())
	} else {
		command := os.Args[2]
		selectedProcessor.Process(command, os.Args[3:])
	}
}

func Help() {
	names := []string{}
	for name := range PROCESSOR_MAP {
		names = append(names, name)
	}

	util.PrintWhite("Valid Processors: " + strings.Join(names, ", "))
}
