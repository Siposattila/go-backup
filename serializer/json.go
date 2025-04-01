package serializer

import "encoding/json"

type jsonSerializer struct{}

func (js *jsonSerializer) Deserialize(data any) ([]byte, error) {
	buffer, marshalError := json.Marshal(data)
	if marshalError != nil {
		return nil, marshalError
	}

	return buffer, nil
}

func (js *jsonSerializer) Serialize(data []byte, structType any) error {
	unMarshalError := json.Unmarshal(data, structType)
	if unMarshalError != nil {
		return unMarshalError
	}

	return nil
}
