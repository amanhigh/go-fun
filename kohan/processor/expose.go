package processor

import (
	"flag"
	"fmt"
	"github.com/amanhigh/go-fun/command"
)

type ExposeProcessor struct {
	Processor
}

func (self *ExposeProcessor) Process(commandName string) (bool) {
	var e error
	flagSet := flag.NewFlagSet(commandName, flag.ExitOnError)

	switch commandName {
	case "pssh":
		e = self.psshHandler(flagSet)
	case "getVersion":
		e = self.getVersionHandler(flagSet)
	default:
		fmt.Println(self.Help())
	}

	if e != nil {
		fmt.Println(e.Error())
		flagSet.Usage()
		return false
	}
	return true
}
func (self *ExposeProcessor) getVersionHandler(flagSet *flag.FlagSet) error {
	pkg := flagSet.String("pkg", "", "Package Name")
	host := flagSet.String("host", "", "Host For Fetching Version")
	versionType := flagSet.String("type", "", "Type dpkg/latest for Version")
	e := flagSet.Parse(self.Args)
	commander.GetVersion(*pkg, *host, *versionType)
	return e
}

func (self *ExposeProcessor) psshHandler(flagSet *flag.FlagSet) error {
	cmd := flagSet.String("cmd", "", "Command To Run")
	cluster := flagSet.String("cl", "", "Cluster To Run On")
	parallelism := flagSet.Int("p", 50, "Parallelism")
	psshType := flagSet.String("t", "fast", "fast/display/slow")
	e := flagSet.Parse(self.Args)
	selectedPssh := getPsshFromType(*psshType)
	selectedPssh.Run(*cmd, *cluster, *parallelism)
	return e
}

func getPsshFromType(psshType string) commander.Pssh {
	var selectedPssh commander.Pssh
	switch psshType {
	case "fast":
		selectedPssh = commander.FastPssh
		break
	case "slow":
		selectedPssh = commander.SlowPssh
	case "display":
		selectedPssh = commander.DisplayPssh

	}
	commander.PrintYellow(fmt.Sprintf("Using %v PSSH", psshType))
	return selectedPssh
}

func (self *ExposeProcessor) Help() string {
	return `Commands: pssh,getVersion`
}
