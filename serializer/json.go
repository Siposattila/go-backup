package serializer

import (
	"encoding/json"
	"log"
)

type jsonSerializer struct{}

func (js *jsonSerializer) Deserialize(data any) []byte {
	buffer, marshalError := json.Marshal(data)
	if marshalError != nil {
		log.Fatal("Unable to desirialize data: " + marshalError.Error())
	}

	return buffer
}

func (js *jsonSerializer) Serialize(data []byte, structType any) {
	unMarshalError := json.Unmarshal(data, structType)
	if unMarshalError != nil {
		log.Fatal("Unable to serialize data: " + unMarshalError.Error())
	}
}
