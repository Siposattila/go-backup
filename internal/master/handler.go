package master

import (
	"context"
	"errors"
	"net/http"

	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/Siposattila/gobkup/internal/request"
	"github.com/Siposattila/gobkup/internal/serializer"
	"github.com/quic-go/webtransport-go"
)

func (master *Master) setupEndpoint() {
	http.HandleFunc("/master", func(writer http.ResponseWriter, request *http.Request) {
		connection, error := master.Server.Upgrade(writer, request)
		if error != nil {
			console.Error("Upgrading failed: " + error.Error())
			writer.WriteHeader(500)

			return
		}

		stream, streamError := connection.AcceptStream(context.Background())
		if streamError != nil {
			console.Error("There was an error during accepting the stream: " + streamError.Error())
		}

		go master.handleStream(stream)
	})
}

func (master *Master) handleStream(clientStream webtransport.Stream) {
	for {
		buffer, readError := master.readFromStream(clientStream)
		if readError != nil {
			break
		}

		var incomingRequest request.NodeRequest
		serializer.Json.Serialize(buffer, &incomingRequest)
		authError := master.authClient(clientStream, incomingRequest.NodeId, incomingRequest.Token)
		if authError != nil {
			clientStream.Close()
			break
		}

		switch incomingRequest.Id {
		case request.REQUEST_ID_CONFIG:
			console.Normal(incomingRequest.NodeId + " sent a request for the config.")
			nodeConfig := config.LoadNodeConfig(incomingRequest.NodeId)
			master.writeToStream(clientStream, master.makeResponse(request.REQUEST_ID_CONFIG, string(serializer.Json.Deserialize(nodeConfig))))
			console.Normal("Config sent to " + incomingRequest.NodeId)
			break
		}
	}
}

func (master *Master) authClient(stream webtransport.Stream, nodeId string, requestToken string) error {
	token, ok := master.Config.Nodes[nodeId]
	if master.Config.RegisterNodeIfKnown && !ok {
		master.AddNode(nodeId)
		master.writeToStream(stream, master.makeResponse(request.REQUEST_ID_NODE_REGISTERED, "TOKEN"))

		return errors.New("NEW")
	}

	if ok && requestToken != token {
		master.writeToStream(stream, master.makeResponse(request.REQUEST_ID_AUTH_ERROR, "TOKEN"))

		return errors.New("AUTH")
	}

	return nil
}

func (master *Master) makeResponse(id int, data string) request.MasterResponse {
	return request.MasterResponse{
		Id:   id,
		Data: data,
	}
}

func (master *Master) writeToStream(stream webtransport.Stream, data any) {
	_, writeError := stream.Write(serializer.Json.Deserialize(data))
	if writeError != nil {
		console.Error("Error during write to node: " + writeError.Error())
	}
}

func (master *Master) readFromStream(stream webtransport.Stream) ([]byte, error) {
	buffer := make([]byte, 1024)
	n, readError := stream.Read(buffer)
	if readError != nil {
		console.Error("Error during reading from node: " + readError.Error())
		return nil, readError
	}

	return buffer[:n], nil
}
