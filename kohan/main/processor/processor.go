package processor

type Processor struct {
}

type ProcessorI interface {
	Process(commandName string) (bool)
	Help() (string)
}
