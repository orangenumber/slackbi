package main

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/rotatew"
	"github.com/gonyyi/graceful"
	"github.com/orangenumber/slackbi"
	"os"
)

func main() {
	// BASIC THING: crate a logger, and writer, graceful
	log := alog.New(os.Stderr)
	r, err := rotatew.New("./log/sbi-shorty{-2006-0102}.log", rotatew.KB*128, rotatew.O_DEFAULT)
	log.IfFatal(err)
	log.SetOutput(r)
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
