package main

import (
	"fmt"
	"os"
	"github.com/amanhigh/go-fun/kohan/main/processor"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: kohan <Processor Name> <Command Name>/help")
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
	return map[string]processor.ProcessorI{
		"ssh": &processor.SshProcessor{},
	}
}
