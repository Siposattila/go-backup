package serializer

type Serializer interface {
	Deserialize(data any) []byte
	Serialize(data []byte, structType any)
}

var Json Serializer = &jsonSerializer{}
