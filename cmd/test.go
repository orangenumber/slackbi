package main

import sbi "github.com/gonyyi/slackbi"

func main() {
	c, err := sbi.ReadConfig("./config.json")
	if err != nil {
		println(err.Error())
	}
	b := sbi.New(*c)
	if err := b.Run(); err!=nil {
		println(err.Error())
	}
	c.Save("./config.json")

}
