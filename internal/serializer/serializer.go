package serializer

import (
	"encoding/json"

	"github.com/Siposattila/gobkup/internal/console"
)

// This only knows json!
func Deserialize(data any) []byte {
	var buffer, marshalError = json.Marshal(data)
	if marshalError != nil {
		console.Fatal("Unable to desirialize data: " + marshalError.Error())
	}

	return buffer
}

// This only knows json!
func Serialize(data []byte, structType any) {
	var unMarshalError = json.Unmarshal(data, structType)
	if unMarshalError != nil {
		console.Fatal("Unable to serialize data: " + unMarshalError.Error())
	}

	return
}
