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

	"github.com/Siposattila/go-backup/generatedproto"
	"github.com/Siposattila/go-backup/log"
	"github.com/Siposattila/go-backup/request"
)

const (
	CHUNK_SIZE     = 10 << 10 // 10KB
	CHUNK_TEMP_DIR = "./chunk_temp"
	CHUNK_NAME     = "%s.part%d"
)

//type Info struct {
//	Name string `json:"name"`
//	Size int    `json:"size"`
//}
//
//type Chunk struct {
//	Name      string `json:"name"`
//	ChunkName string `json:"chunkName"`
//	Data      []byte `json:"data"`
//	Size      int    `json:"size"`
//}

func NewInfo(name string, size int32) *generatedproto.Info {
	return &generatedproto.Info{
		Name: name,
		Size: size,
	}
}

func NewChunk(name string, chunkName string, data []byte) *generatedproto.Chunk {
	return &generatedproto.Chunk{
		Name:      name,
		ChunkName: chunkName,
		Data:      data,
		Size:      CHUNK_SIZE,
	}
}

func chunkFile(name string) (int, error) {
	if _, err := os.Stat(CHUNK_TEMP_DIR); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(CHUNK_TEMP_DIR, 0777); err != nil {
			log.GetLogger().Fatal("Stat failed during chunkFile: ", err.Error())
		}
	}

	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}

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
			if err := partFile.Close(); err != nil {
				log.GetLogger().Error("Writing part file failed: ", err.Error())
			}
			return 0, writeErr
		}

		if err := partFile.Close(); err != nil {
			log.GetLogger().Error("Closing part file failed: ", err.Error())
		}
		partNum++
	}

	if err := file.Close(); err != nil {
		log.GetLogger().Error("Closing file failed: ", err.Error())
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
			log.GetLogger().Fatal("Stat failed for new backup path: ", err.Error())
		}

		backupInfo := NewInfo(backupName, int32(info.Size()))
		if _, err := request.Write(c.Stream, request.NewRequest(
			c.Config.ClientId,
			generatedproto.RequestType_ID_BACKUP_START,
			backupInfo)); err != nil {
			log.GetLogger().Fatal("Writing backup info failed: ", err.Error())
		}
		// TODO: should investigate if we need this or not
		// time.Sleep(10 * time.Millisecond) // looks like this is necessary because it writes too fast

		n, err := chunkFile(newBackupPath)
		if err != nil {
			log.GetLogger().Fatal("Chunk file failed: ", err.Error())
		}

		for partNum := range n {
			partFile, err := os.Open(path.Join(CHUNK_TEMP_DIR, fmt.Sprintf(CHUNK_NAME, backupName, partNum)))
			if err != nil {
				log.GetLogger().Fatal("Opening part file failed: ", err.Error())
			}

			data := make([]byte, CHUNK_SIZE)
			if _, err := partFile.Read(data); err != nil {
				log.GetLogger().Fatal("Reading part file failed: ", err.Error())
			}
			if _, err := request.Write(
				c.Stream,
				request.NewRequest(
					c.Config.ClientId,
					generatedproto.RequestType_ID_BACKUP_CHUNK,
					NewChunk(
						backupName,
						strings.Split(partFile.Name(), "/")[1], data))); err != nil {
				log.GetLogger().Fatal("Writing chunk to stream failed: ", err.Error())
			}

			if err := partFile.Close(); err != nil {
				log.GetLogger().Error("Closing part file failed: ", err.Error())
			}
			time.Sleep(10 * time.Millisecond) // looks like this is necessary because it writes too fast
		}

		if _, err := request.Write(c.Stream, request.NewRequest(c.Config.ClientId, generatedproto.RequestType_ID_BACKUP_END, backupInfo)); err != nil {
			log.GetLogger().Fatal("Writing backup end to stream failed: ", err.Error())
		}
	}
}
