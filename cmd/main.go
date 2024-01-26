package main

import (
	"flag"

	"github.com/Siposattila/gobkup/internal/console"
)

func main() {
	//var isMaster = flag.Bool("master", false, "This flag will start gobkup in master mode which means this is going to be the master server.")
	//var isNode = flag.Bool("node", false, "This flag will start gobkup in node mode which means you should also configure the master endpoint with the --endpoint flag.")
	//var endpoint = flag.String("endpoint", "", "This flag will set the master endpoint.")
	//var isGenerate = flag.Bool("generate", false, "This flag will generate the master config.")
	//var nodeName = flag.String("add-node", "", "This flag will add a node. The node id you add should be the name of the server. --add-node <NodeId>")

	flag.Parse()

	if isFlagPassed("master") {
		// TODO: Should load the master config
		console.Debug("master")
	}

	if isFlagPassed("node") && isFlagPassed("endpoint") {
		// TODO: Connect to the master endpoint get the config for the node
		console.Debug("node and endpoint")
	}

	if isFlagPassed("generate") {
		// TODO: Master config generation
		console.Debug("generate")
	}

	if isFlagPassed("add-node") {
		// TODO: Add node to master
		console.Debug("add node")
	}

	if !isFlagPassed("master") && !isFlagPassed("node") && !isFlagPassed("generate") && !isFlagPassed("add-node") {
		console.White("No flags were provided. Please use the -h or --help flag for help.")
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
