package command

import "time"

var (
	cluster string
	pkgName string
	command string
	marker  string
	tyype   string

	parallelism = -1
	index       = -1
	endIndex    = -1
	year        = -1
	cutOff      = -1
	count       = -1

	//Auto
	wait = 5 * time.Second
	idle = 60 * time.Second

	enable  = false
	verbose = false
)
