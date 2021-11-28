package main

import (
	"github.com/liuguangw/billing_go/cmd"
	"log"
	"os"
)

func main() {
	appCommand := cmd.AppCommand()
	if err := appCommand.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
