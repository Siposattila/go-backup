package request

import (
	"github.com/Siposattila/go-backup/serializer"
	"github.com/quic-go/webtransport-go"
)

const (
	ID_CONFIG                 = 10010
	ID_BACKUP_START           = 10020
	ID_BACKUP_CHUNK           = 10030
	ID_BACKUP_CHUNK_PROCESSED = 10040
	ID_BACKUP_END             = 10050
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
	buffer := make([]byte, 50<<10) // 50 KB
	totalRead := 0
	isFull := false
	for !isFull {
		n, readError := stream.Read(buffer[totalRead:])
		if readError != nil {
			return 0, readError
		}

		totalRead += n
		serializerError := serializer.Json.Serialize(buffer[:totalRead], data)
		if serializerError == nil {
			isFull = true
		}
	}

	return totalRead, nil
}
