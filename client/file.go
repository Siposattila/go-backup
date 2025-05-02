package client

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Siposattila/gobkup/log"
	"github.com/Siposattila/gobkup/request"
)

const (
	CHUNK_SIZE     = 10 << 10 // 10KB
	CHUNK_TEMP_DIR = "./chunk_temp"
	CHUNK_NAME     = "%s.part%d"
)

type Info struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

type Chunk struct {
	Name      string `json:"name"`
	ChunkName string `json:"chunkName"`
	Data      []byte `json:"data"`
	Size      int    `json:"size"`
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
	if _, err := os.Stat(CHUNK_TEMP_DIR); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(CHUNK_TEMP_DIR, 0777)
	}

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

func (c *client) handNewBackup() {
	for {
		newBackupPath := <-c.newBackupPathChannel
		helper := strings.Split(newBackupPath, "/")
		backupName := helper[len(helper)-1]

		info, err := os.Stat(newBackupPath)
		if err != nil {
			log.GetLogger().Fatal(err.Error())
		}

		backupInfo := NewInfo(backupName, int(info.Size()))
		request.Write(c.Stream, request.NewRequest(c.Config.ClientId, request.ID_BACKUP_START, backupInfo))
		time.Sleep(10 * time.Millisecond) // looks like this is necessary because it writes too fast

		n, err := chunkFile(newBackupPath)
		if err != nil {
			log.GetLogger().Fatal(err.Error())
		}

		for partNum := range n {
			partFile, err := os.Open(path.Join(CHUNK_TEMP_DIR, fmt.Sprintf(CHUNK_NAME, backupName, partNum)))
			if err != nil {
				log.GetLogger().Fatal(err.Error())
			}

			data := make([]byte, CHUNK_SIZE)
			partFile.Read(data)
			request.Write(c.Stream, request.NewRequest(c.Config.ClientId, request.ID_BACKUP_CHUNK, NewChunk(backupName, strings.Split(partFile.Name(), "/")[1], data)))

			partFile.Close()
			time.Sleep(10 * time.Millisecond) // looks like this is necessary because it writes too fast
		}

		request.Write(c.Stream, request.NewRequest(c.Config.ClientId, request.ID_BACKUP_END, backupInfo))
	}
}
