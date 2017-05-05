package processor

type Processor struct {
	Args []string
}

type ProcessorI interface {
	Process(commandName string) (bool)
	Help() (string)
}
