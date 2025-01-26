package main

import (
	"flag"

	"github.com/Siposattila/gobkup/dealer"
	"github.com/Siposattila/gobkup/log"
)

func main() {
	server := flag.Bool("server", false, "This flag will start gobkup in server mode.")
	client := flag.Bool("client", false, "This flag will start gobkup in client mode.")
	clientName := flag.String("add-client", "", "With this flag you can add a client. The client id should be the name of the server. --add-client <ClientId>")

	flag.Parse()

	log.GetLogger().Debug(*server, *client, *clientName, *clientName != "")

	if *server || *client {
		dealer.Run(*server, *client)
	} else {
		if *clientName != "" {
			// TODO: logic for adding a client
		}
	}

	log.GetLogger().Fatal("Expected atleast one valid flag!")
}
