package request

import (
	"github.com/Siposattila/gobkup/serializer"
	"github.com/quic-go/webtransport-go"
)

const (
	REQUEST_ID_CONFIG          = 10010
	REQUEST_ID_NODE_REGISTERED = 10020
	REQUEST_ID_AUTH_ERROR      = 10030
	REQUEST_ID_KEEPALIVE       = 10040
)

type Response struct {
	Id   int    `json:"id"`
	Data string `json:"data"`
}

type Request struct {
	Id       int    `json:"id"`
	Data     string `json:"data"`
	DealerId string `json:"dealerId"`
	// Token  string `json:"token"`
}

func NewResponse(id int, data string) *Response {
	return &Response{
		Id:   id,
		Data: data,
	}
}

func NewRequest(id int, data string) *Request {
	return &Request{
		Id:   id,
		Data: data,
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

func Read[T Request | Response](stream webtransport.Stream, data T) (int, error) {
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
