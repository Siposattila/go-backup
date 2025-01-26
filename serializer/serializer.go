package serializer

type Serializer interface {
	Deserialize(data any) ([]byte, error)
	Serialize(data []byte, structType any) error
}

var Json Serializer = &jsonSerializer{}
