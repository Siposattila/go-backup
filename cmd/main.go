package main

import (
	"flag"

	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/Siposattila/gobkup/internal/node"
)

func main() {
	flag.Bool("master", false, "This flag will start gobkup in master mode which means this is going to be the master server.")
	flag.Bool("node", false, "This flag will start gobkup in node mode which means you should also configure the master endpoint with the --endpoint flag.")
	flag.Bool("generate", false, "This flag will generate the master config.")
	flag.Bool("debug", false, "This flag will set the debug mode to true. If debug mode is set then the server will not use tls!")
	var endpoint = flag.String("endpoint", "", "This flag will set the master endpoint.")
	var nodeName = flag.String("add-node", "", "This flag will add a node. The node id you add should be the name of the server. --add-node <NodeId>")

	flag.Parse()

	if isFlagPassed("master") {
		config.LoadConfig("master")
        config.Master.Debug = isFlagPassed("debug")
		// TODO: should start the server
	}

	if isFlagPassed("node") && isFlagPassed("endpoint") {
		config.LoadConfig("node")
        config.Node.Debug = isFlagPassed("debug")
		// TODO: Connect to the master endpoint get the config for the node
		console.Debug(*endpoint)
		console.Debug("node and endpoint")
	}

	if isFlagPassed("generate") {
		config.GenerateMasterConfig()
	}

	if isFlagPassed("add-node") {
		if *nodeName == "" {
			console.Fatal("You must provide the name(or id) of the node.")
		}
		node.AddNode(*nodeName)
	}
}

func isFlagPassed(name string) bool {
	var found = false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})

	return found
}
