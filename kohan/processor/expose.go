package processor

import (
	"flag"
	"github.com/amanhigh/go-fun/kohan/commander/components"
)

type ExposeProcessor struct {
}

func (self *ExposeProcessor) GetArgedHandlers() (map[string]HandleFunc) {
	return map[string]HandleFunc{
		"getVersion": self.getVersionHandler,
		"printf":     self.handlePrintf,
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

func (self *ExposeProcessor) handlePrintf(flagSet *flag.FlagSet, args []string) error {
	templateFile := flagSet.String("c", "", "Template File")
	paramFile := flagSet.String("p", "", "Params File")
	marker := flagSet.String("m", "#", "Marker")
	e := flagSet.Parse(args)
	components.Printf(*templateFile, *paramFile, *marker)
	return e
}
