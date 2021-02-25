package goc2

import (
	"flag"
	//"fmt"
	//"os"
	"goc2/pkg/cli"
	"goc2/web"
)

var (
	cliPtr   bool
	webPtr   bool
	listPtr  bool
	agentPtr string
	c2Ptr    string
	portPtr  string
)

//Start RedMap
func Start() {
	//flags
	flag.BoolVar(&cliPtr, "cli", false, "start cli")
	flag.BoolVar(&webPtr, "web", false, "Start Web Server")
	flag.BoolVar(&listPtr, "list", false, "List Connected Agents")
	flag.StringVar(&agentPtr, "agent", "", "create payload")
	flag.StringVar(&portPtr, "port", "8005", "Listen Port")
	flag.StringVar(&c2Ptr, "c2", "", "connect to c2")
	flag.Parse()

	if listPtr == true {
		cli.ListAgents(c2Ptr)
	}

	if webPtr == true {
		web.Start(portPtr)
	}

	if cliPtr == true {
		cli.Start(c2Ptr)
	}

}
