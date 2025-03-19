package main

import (
	"flag"
	"frappuccino/app"
	"frappuccino/tools"
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

	app.StartServer()
}
