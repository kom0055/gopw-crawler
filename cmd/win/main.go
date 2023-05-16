package main

import (
	"gopw-crawler/cmd"
	"os"
)

func main() {
	c := cmd.NewDailyPowerQueryCmd()
	err := c.Execute()
	if err != nil {
		os.Exit(1)
	}
}
