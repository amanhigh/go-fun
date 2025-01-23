package command

import "time"

var (
	cluster string
	command string
	tyype   string

	parallelism = -1
	index       = -1
	endIndex    = -1

	//Auto
	wait = time.Minute

	// Darius
	makeFileDir    = "/home/aman/Projects/go-fun/Kubernetes/"
	tmpServiceFile = "/tmp/k8-svc.txt"
)
