package main

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/arotw"
	"github.com/gonyyi/graceful"
	"github.com/orangenumber/slackbi"
	"os"
)

func main() {
	// BASIC THING: crate a logger, and writer, graceful
	log := alog.New(os.Stderr)
	r, err := arotw.New("./log/sbi-shorty{-2006-0102}.log", func(rw *arotw.Writer) {
		rw.SetKeepLogs(5)
		rw.SetMaxSize(arotw.MB * 10)
	})

	log.SetOutput(r).Do(alog.DoColor()).SetLevel(alog.Ltrace)
	log.IfFatal(err)

	graceful.New(func() {
		log.Fatal("received a shutdown signal")
	})

	// SBI CONFIG
	sbiconf, err := slackbi.ReadConfig("config.json")
	log.IfError(err)

	sbi, err := slackbi.New(sbiconf, log)
	log.IfFatal(err)

	log.IfFatal(sbi.Run())
}
