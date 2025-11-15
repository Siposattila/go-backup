package request

import (
	"github.com/Siposattila/go-backup/generatedproto"
	"github.com/quic-go/webtransport-go"
	"google.golang.org/protobuf/proto"
)

func NewResponse(id generatedproto.RequestType, data proto.Message) *generatedproto.Response {
	d, _ := proto.Marshal(data)

	return &generatedproto.Response{
		Id:   id,
		Data: d,
	}
}

func NewRequest(clientId string, id generatedproto.RequestType, data proto.Message) *generatedproto.Request {
	d, _ := proto.Marshal(data)

	return &generatedproto.Request{
		Id:       id,
		Data:     d,
		ClientId: clientId,
	}
}

func Write(stream webtransport.Stream, data proto.Message) (int, error) {
	dataToWrite, serializerError := proto.Marshal(data)
	if serializerError != nil {
		return 0, serializerError
	}

	n, writeError := stream.Write(dataToWrite)
	if writeError != nil {
		return 0, writeError
	}

	return n, nil
}

func Read(stream webtransport.Stream, data proto.Message) (int, error) {
	buffer := make([]byte, 500<<10) // 500 KB
	n, readError := stream.Read(buffer)
	if readError != nil {
		return 0, readError
	}

	if unmarshalError := proto.Unmarshal(buffer, data); unmarshalError != nil {
		return 0, unmarshalError
	}

	return n, nil
}
