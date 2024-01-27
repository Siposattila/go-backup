package master

import (
	"errors"

	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/Siposattila/gobkup/internal/node"
	"github.com/Siposattila/gobkup/internal/request"
	"github.com/Siposattila/gobkup/internal/serializer"
	"github.com/quic-go/webtransport-go"
)

func handleStream(clientStream webtransport.Stream) {
	for {
		var buffer, readError = readFromStream(clientStream)
		if readError != nil {
			break
		}

		var incomingRequest request.NodeRequest
		serializer.Serialize(buffer, &incomingRequest)
		var authError = authClient(clientStream, incomingRequest.NodeId, incomingRequest.Token)
        if authError != nil {
            clientStream.Close()
            break
        }

		switch incomingRequest.Id {
		case "Config":
            console.Normal(incomingRequest.NodeId + " sent a request for the config.")
			config.LoadConfig("node", incomingRequest.NodeId)
			writeToStream(clientStream, makeResponse("Config", string(serializer.Deserialize(config.Node))))
            console.Normal("Config sent to " + incomingRequest.NodeId)
			break
		}
	}

	return
}

func authClient(stream webtransport.Stream, nodeId string, requestToken string) error {
	var token, ok = config.Master.Nodes[nodeId]
	if config.Master.RegisterNodeIfKnown && !ok {
		node.AddNode(nodeId)
		writeToStream(stream, makeResponse("NewNode", "TOKEN"))

        return errors.New("NEW")
	}

	if ok && requestToken != token {
		writeToStream(stream, makeResponse("NotValid", "ERROR"))

        return errors.New("INVALID")
	}

    return nil
}

func makeResponse(id string, data string) request.MasterResponse {
	return request.MasterResponse{
		Id:   id,
		Data: data,
	}
}

func writeToStream(stream webtransport.Stream, data any) {
	var _, writeError = stream.Write(serializer.Deserialize(data))
	if writeError != nil {
		console.Error("Error during write to node: " + writeError.Error())
	}

	return
}

func readFromStream(stream webtransport.Stream) ([]byte, error) {
	var buffer = make([]byte, 1024)
	var n, readError = stream.Read(buffer)
	if readError != nil {
		console.Error("Error during reading from node: " + readError.Error())
		return nil, readError
	}

	return buffer[:n], nil
}
