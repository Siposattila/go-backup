package client

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

const (
	CHUNK_SIZE     = 1 * 1024 * 1024 // 1MB
	CHUNK_TEMP_DIR = "./chunk_temp"
	CHUNK_NAME     = "%s.part%d"
)

type Info struct {
	Name string
	Size int
}

type Chunk struct {
	Name      string
	ChunkName string
	Data      []byte
	Size      int
}

func NewInfo(name string, size int) *Info {
	return &Info{
		Name: name,
		Size: size,
	}
}

func NewChunk(name string, chunkName string, data []byte) *Chunk {
	return &Chunk{
		Name:      name,
		ChunkName: chunkName,
		Data:      data,
		Size:      CHUNK_SIZE,
	}
}

func chunkFile(name string) (int, error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buffer := make([]byte, CHUNK_SIZE)
	baseName := filepath.Base(name)

	partNum := 0
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			break
		}

		partFileName := fmt.Sprintf(CHUNK_NAME, baseName, partNum)
		partFile, err := os.Create(path.Join(CHUNK_TEMP_DIR, partFileName))
		if err != nil {
			return 0, err
		}

		_, writeErr := partFile.Write(buffer[:n])
		if writeErr != nil {
			partFile.Close()
			return 0, writeErr
		}

		partFile.Close()
		partNum++
	}

	return partNum, nil
}
