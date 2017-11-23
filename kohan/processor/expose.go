package processor

import (
	"flag"
	"fmt"
	"github.com/amanhigh/go-fun/kohan/commander"
	"github.com/amanhigh/go-fun/kohan/commander/components"
	. "github.com/amanhigh/go-fun/kohan/commander/tools"
	"github.com/amanhigh/go-fun/util"
)

type ExposeProcessor struct {
}

func (self *ExposeProcessor) GetArgedHandlers() (map[string]HandleFunc) {
	return map[string]HandleFunc{
		"pssh":         self.psshHandler,
		"getVersion":   self.getVersionHandler,
		"indexedIp":    self.handleIndexedIp,
		"printf":       self.handlePrintf,
		"versionCheck": self.versionCheckHandler,
		"verifyStatus": self.verifyStatusHandler,
		"debugControl": self.debugControlHandler,
	}
}

func (self *ExposeProcessor) GetNonArgedHandlers() (map[string]DirectFunc) {
	return map[string]DirectFunc{}
}

func (self *ExposeProcessor) getVersionHandler(flagSet *flag.FlagSet, args []string) error {
	pkg := flagSet.String("pkg", "", "Package Name")
	host := flagSet.String("host", "", "Host For Fetching Version")
	versionType := flagSet.String("type", "", "Type dpkg/latest for Version")
	comment := flagSet.String("c", "N/A", "Comment for this release")
	e := flagSet.Parse(args)
	components.GetVersion(*pkg, *host, *versionType, *comment)
	return e
}

func (self *ExposeProcessor) handleIndexedIp(flagSet *flag.FlagSet, args []string) error {
	cluster := flagSet.String("cl", "", "Cluster Name")
	index := flagSet.Int("i", -1, "Index of Ip")
	e := flagSet.Parse(args)
	IndexedIp(*cluster, *index)
	return e
}

func (self *ExposeProcessor) versionCheckHandler(flagSet *flag.FlagSet, args []string) error {
	pkg := flagSet.String("pkg", "", "CSV List of Package Names")
	cluster := flagSet.String("cl", "", "Cluster To Run On")
	e := flagSet.Parse(args)
	components.VersionCheck(*pkg, *cluster)
	return e
}

func (self *ExposeProcessor) verifyStatusHandler(flagSet *flag.FlagSet, args []string) error {
	cmd := flagSet.String("cmd", "", "Status Check Command")
	cluster := flagSet.String("cl", "", "Cluster To Run On")
	e := flagSet.Parse(args)
	components.VerifyStatus(*cmd, *cluster)
	return e
}

func (self *ExposeProcessor) handlePrintf(flagSet *flag.FlagSet, args []string) error {
	template := flagSet.String("c", "", "Command Template")
	params := flagSet.String("p", "", "Params to Merge")
	marker := flagSet.String("m", "#", "Marker")
	e := flagSet.Parse(args)
	components.Printf(*template, *params, *marker)
	return e
}

func (self *ExposeProcessor) psshHandler(flagSet *flag.FlagSet, args []string) error {
	cmd := flagSet.String("cmd", "", "Command To Run")
	cluster := flagSet.String("cl", "", "Cluster To Run On")
	psshType := flagSet.String("t", "fast", "fast/display/slow")
	parallelism := flagSet.Int("p", commander.DEFAULT_PARALELISM, "Parallelism")
	start := flagSet.Int("s", -1, "Start Index (Starting from 1)")
	end := flagSet.Int("e", -1, "End Index")
	e := flagSet.Parse(args)
	selectedPssh := getPsshFromType(*psshType)
	selectedPssh.RunRange(*cmd, *cluster, *parallelism, false, *start, *end)
	return e
}

func (self *ExposeProcessor) debugControlHandler(flagSet *flag.FlagSet, args []string) error {
	f := flagSet.Bool("f", false, "Enable Disable Flag true/false")
	e := flagSet.Parse(args)
	commander.DebugControl(*f)
	return e
}

func getPsshFromType(psshType string) Pssh {
	var selectedPssh Pssh
	switch psshType {
	case "fast":
		selectedPssh = FastPssh
		break
	case "slow":
		selectedPssh = SlowPssh
	case "display":
		selectedPssh = DisplayPssh

	}
	util.PrintYellow(fmt.Sprintf("Using %v PSSH", psshType))
	return selectedPssh
}
