package main

import sbi "github.com/orangenumber/slackbi"

func main() {
	c, err := sbi.ReadConfig("./config.json")
	if err != nil {
		println(err.Error())
		return
	}
	b := sbi.New(*c)
	if err := b.Run(); err != nil {
		println(err.Error())
	}

	b.Run()

}
