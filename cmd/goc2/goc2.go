package goc2

import (
	"flag"
	//"fmt"
	//"os"
	"goc2/web"
	"goc2/pkg/cli"
)

var (
	cliPtr 		bool
	webPtr      bool
)

//Start RedMap
func Start() {
	//flags
	flag.BoolVar(&cliPtr, "cli", false, "run email check")
	flag.BoolVar(&webPtr, "web", false, "Start Web Server")
	flag.Parse()

	if webPtr == true {
		web.Start()
	}

	if cliPtr == true {
		cli.Start()
	}

}
