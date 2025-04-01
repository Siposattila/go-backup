package request

import (
	"github.com/Siposattila/gobkup/serializer"
	"github.com/quic-go/webtransport-go"
)

const (
	REQUEST_ID_CONFIG = 10010
)

type Response struct {
	Id   int    `json:"id"`
	Data string `json:"data"`
}

type Request struct {
	Id       int    `json:"id"`
	Data     string `json:"data"`
	ClientId string `json:"clientId"`
}

func NewResponse(id int, data any) *Response {
	d, _ := serializer.Json.Deserialize(data)

	return &Response{
		Id:   id,
		Data: string(d),
	}
}

func NewRequest(clientId string, id int, data any) *Request {
	d, _ := serializer.Json.Deserialize(data)

	return &Request{
		Id:       id,
		Data:     string(d),
		ClientId: clientId,
	}
}

func Write[T *Request | *Response](stream webtransport.Stream, data T) (int, error) {
	dataToWrite, serializerError := serializer.Json.Deserialize(data)
	if serializerError != nil {
		return 0, serializerError
	}

	n, writeError := stream.Write(dataToWrite)
	if writeError != nil {
		return 0, writeError
	}

	return n, nil
}

func Read[T *Request | *Response](stream webtransport.Stream, data T) (int, error) {
	buffer := make([]byte, 1024)
	n, readError := stream.Read(buffer)
	if readError != nil {
		return 0, readError
	}

	serializerError := serializer.Json.Serialize(buffer[:n], data)
	if serializerError != nil {
		return 0, serializerError
	}

	return n, nil
}
