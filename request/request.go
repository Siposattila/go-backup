package request

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/quic-go/webtransport-go"
	"google.golang.org/protobuf/proto"
)

func Write(stream webtransport.Stream, msg proto.Message) (int, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return 0, fmt.Errorf("marshal: %w", err)
	}

	// Write length prefix
	length := uint32(len(data))
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, length)

	if _, err := stream.Write(lenBuf); err != nil {
		return 0, fmt.Errorf("write length: %w", err)
	}

	n, err := stream.Write(data)
	if err != nil {
		return 0, fmt.Errorf("write data: %w", err)
	}

	return n, nil
}

func Read(stream webtransport.Stream, msg proto.Message) (int, error) {
	// Read the length prefix (4 bytes)
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(stream, lenBuf); err != nil {
		return 0, fmt.Errorf("read length: %w", err)
	}

	length := binary.BigEndian.Uint32(lenBuf)
	if length == 0 {
		return 0, fmt.Errorf("invalid zero-length message")
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(stream, data); err != nil {
		return 0, fmt.Errorf("read payload: %w", err)
	}

	if err := proto.Unmarshal(data, msg); err != nil {
		return 0, fmt.Errorf("unmarshal: %w", err)
	}

	return int(length), nil
}
