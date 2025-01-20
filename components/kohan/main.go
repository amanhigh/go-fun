package main

import (
	"github.com/amanhigh/go-fun/components/kohan/command"
	"github.com/amanhigh/go-fun/components/kohan/core"
)

func main() {
	core.SetupKohanInjector()

	command.Execute()
}
