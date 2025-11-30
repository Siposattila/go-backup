package client

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Siposattila/go-backup/log"
	"github.com/Siposattila/go-backup/proto"
	"github.com/Siposattila/go-backup/request"
)

const (
	CHUNK_SIZE     = 10 << 10 // 10 kilobytes
	CHUNK_TEMP_DIR = "./chunk_temp"
	CHUNK_NAME     = "%s.part%d"
)

func chunkFile(name string) (int, error) {
	if _, err := os.Stat(CHUNK_TEMP_DIR); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(CHUNK_TEMP_DIR, 0777); err != nil {
			return 0, fmt.Errorf("failed to mkdir chunk temp directory: %s", err.Error())
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
				return 0, fmt.Errorf("closing part file failed after failing to write to it: %s", err.Error())
			}
			return 0, fmt.Errorf("writing part file failed: %s", writeErr.Error())
		}

		if err := partFile.Close(); err != nil {
			return 0, fmt.Errorf("closing part file failed: %s", err.Error())
		}
		partNum++
	}

	if err := file.Close(); err != nil {
		return 0, fmt.Errorf("closing file %s failed: %s", name, err.Error())
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

		if _, err := request.Write(c.ServerStream, &proto.Envelope{
			ClientId: c.Config.ClientId,
			Message: &proto.Envelope_BackupStartRequest{
				BackupStartRequest: &proto.BackupStartRequest{
					Name: backupName,
					Size: int32(info.Size()),
				},
			},
		}); err != nil {
			log.GetLogger().Fatal("Writing backup info failed: ", err.Error())
		}
		// TODO: should investigate if we need this or not
		// time.Sleep(10 * time.Millisecond) // Looks like this is necessary because it writes too fast

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
			size, partFileReadError := partFile.Read(data)
			if partFileReadError != nil {
				log.GetLogger().Fatal("Reading part file failed: ", partFileReadError.Error())
			}
			if _, err := request.Write(
				c.ServerStream, &proto.Envelope{
					ClientId: c.Config.ClientId,
					Message: &proto.Envelope_BackupChunkRequest{
						BackupChunkRequest: &proto.BackupChunkRequest{
							Chunk: &proto.BackupChunk{
								Name:      backupName,
								ChunkName: strings.Split(partFile.Name(), "/")[1],
								Data:      data,
								Size:      int32(size),
							},
						},
					},
				}); err != nil {
				log.GetLogger().Fatal("Writing chunk to stream failed: ", err.Error())
			}

			if err := partFile.Close(); err != nil {
				log.GetLogger().Error("Closing part file failed: ", err.Error())
			}
			// TODO: should investigate if we need this or not
			// time.Sleep(10 * time.Millisecond) // Looks like this is necessary because it writes too fast
		}

		if _, err := request.Write(c.ServerStream, &proto.Envelope{
			ClientId: c.Config.ClientId,
			Message: &proto.Envelope_BackupEndRequest{
				BackupEndRequest: &proto.BackupEndRequest{Name: backupName},
			},
		}); err != nil {
			log.GetLogger().Fatal("Writing backup end to stream failed: ", err.Error())
		}
	}
}
