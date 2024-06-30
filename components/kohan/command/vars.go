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
	wait = time.Minute

	enable  = false
	verbose = false

	// Darius
	makeFileDir    = "/home/aman/Projects/go-fun/Kubernetes/"
	tmpServiceFile = "/tmp/k8-svc.txt"
)
