package tools

import "github.com/bitfield/script"

func Screenshot() (err error) {
	_, err = script.Exec("spectacle -mbn").String()
	return
}
