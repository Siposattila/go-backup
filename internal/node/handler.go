package node

import (
	"github.com/Siposattila/gobkup/internal/backup"
	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/Siposattila/gobkup/internal/request"
	"github.com/Siposattila/gobkup/internal/serializer"
)

func handleStream() {
	console.Normal("Trying to request for config from master...")
	writeToStream(makeRequest("Config", "PLEASE"))
	for {
		var response request.MasterResponse
		serializer.Serialize(readFromStream(), &response)

		switch response.Id {
		case "Config": 
			serializer.Serialize([]byte(response.Data), &config.Node)
			console.Success("Got config from master!")
            backup.BackupProcess()
        	writeToStream(makeRequest("KeepAlive", "OK"))
			break
        case "NewNode":
            console.Warning("This node is now registered at the master! The token was generated for this node at the master.")
            return
		case "NotValid":
			console.Fatal("Wrong credentials!")
			break
        case "KeepAlive":
        	writeToStream(makeRequest("KeepAlive", "OK"))
            break
		}
	}
}

func makeRequest(id string, data string) request.NodeRequest {
	return request.NodeRequest{
		Id:     id,
		Data:   data,
		NodeId: getClientName(),
		Token:  config.Node.Token,
	}
}

func writeToStream(data any) {
	var _, writeError = serverStream.Write(serializer.Deserialize(data))
	if writeError != nil {
		console.Error("Error during write to master: " + writeError.Error())
	}
}

func readFromStream() []byte {
	var buffer = make([]byte, 1024)
	var n, readError = serverStream.Read(buffer)
	if readError != nil {
		console.Error("Error during reading from master: " + readError.Error())
	}

	return buffer[:n]
}
