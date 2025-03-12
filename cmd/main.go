package main

import (
	"flag"
	"frappuccino/app"
	"frappuccino/tools"
	"log"
	"os"
)

func main() {
	flag.Parse()
	tools.ParseFlag()

	if *tools.Help {
		tools.HelpFunck()
		flag.Usage()
		os.Exit(0)
	}

	if !tools.CheckDir(*tools.Dir) {
		log.Println("There is no such directory")
		return
	}

	app.StartServer()
}
