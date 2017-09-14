package processor

import (
	"flag"
	"fmt"
	"github.com/amanhigh/go-fun/command"
)

type ExposeProcessor struct {
}

func (self *ExposeProcessor) GetHandleMap() (map[string]HandleFunc) {
	return map[string]HandleFunc{
		"pssh":         self.psshHandler,
		"getVersion":   self.getVersionHandler,
		"indexedIp":    self.handleIndexedIp,
		"versionCheck": self.versionCheckHandler,
		"verifyStatus": self.verifyStatusHandler,
		"debugControl": self.debugControlHandler,
	}
}

func (self *ExposeProcessor) getVersionHandler(flagSet *flag.FlagSet, args []string) error {
	pkg := flagSet.String("pkg", "", "Package Name")
	host := flagSet.String("host", "", "Host For Fetching Version")
	versionType := flagSet.String("type", "", "Type dpkg/latest for Version")
	e := flagSet.Parse(args)
	commander.GetVersion(*pkg, *host, *versionType)
	return e
}

func (self *ExposeProcessor) handleIndexedIp(flagSet *flag.FlagSet, args []string) error {
	cluster := flagSet.String("cl", "", "Cluster Name")
	index := flagSet.Int("i", -1, "Index of Ip")
	e := flagSet.Parse(args)
	commander.IndexedIp(*cluster, *index)
	return e
}

func (self *ExposeProcessor) versionCheckHandler(flagSet *flag.FlagSet, args []string) error {
	pkg := flagSet.String("pkg", "", "CSV List of Package Names")
	cluster := flagSet.String("cl", "", "Cluster To Run On")
	e := flagSet.Parse(args)
	commander.VersionCheck(*pkg, *cluster)
	return e
}

func (self *ExposeProcessor) verifyStatusHandler(flagSet *flag.FlagSet, args []string) error {
	cmd := flagSet.String("cmd", "", "Status Check Command")
	cluster := flagSet.String("cl", "", "Cluster To Run On")
	e := flagSet.Parse(args)
	commander.VerifyStatus(*cmd, *cluster)
	return e
}

func (self *ExposeProcessor) psshHandler(flagSet *flag.FlagSet, args []string) error {
	cmd := flagSet.String("cmd", "", "Command To Run")
	cluster := flagSet.String("cl", "", "Cluster To Run On")
	parallelism := flagSet.Int("p", commander.DEFAULT_PARALELISM, "Parallelism")
	psshType := flagSet.String("t", "fast", "fast/display/slow")
	e := flagSet.Parse(args)
	selectedPssh := getPsshFromType(*psshType)
	selectedPssh.Run(*cmd, *cluster, *parallelism, false)
	return e
}

func (self *ExposeProcessor) debugControlHandler(flagSet *flag.FlagSet, args []string) error {
	f := flagSet.Bool("f", false, "Enable Disable Flag true/false")
	e := flagSet.Parse(args)
	commander.DebugControl(*f)
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
