package main

import (
	"github.com/gonyyi/alog"
	sbi "github.com/orangenumber/slackbi"
	"os"
)

func main() {
	log := alog.New(os.Stderr, "", alog.FDefault|alog.FLevelTrace)

	c, err := sbi.ReadConfig("./config.json")
	if err != nil {
		println(err.Error())
		return
	}
	b, err := sbi.New(c, log)
	if err := b.Run(); err != nil {
		println(err.Error())
	}

	b.Run()

}
