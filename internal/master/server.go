package master

import (
	"net/http"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func SetupAndRunServer() {
	var server = webtransport.Server {
		H3: http3.Server{Addr: ":443"},
	}

    // TODO: implement master
	http.HandleFunc("/master", func(writer http.ResponseWriter, request *http.Request) {
		// var connection, error = server.Upgrade(writer, request)
		//if error != nil {
		//	console.Error("Upgrading failed: " + error.Error())
		//	writer.WriteHeader(500)

		//	return
		//}
		// Handle the connection. Here goes the application logic.
	})

	// server.ListenAndServeTLS(certFile, keyFile)
    server.ListenAndServe()
}
