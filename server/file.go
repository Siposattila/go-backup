package server

import (
	"fmt"
	"os"

	"github.com/Siposattila/go-backup/generatedproto"
	"github.com/Siposattila/go-backup/log"
)

const TEMP_FILE = "%s.temp"

func (s *server) writeChunk(chunk *generatedproto.Chunk) {
	name := fmt.Sprintf("%s/%s", s.Config.BackupPath, fmt.Sprintf(TEMP_FILE, chunk.Name))
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.GetLogger().Error(err)
	}

	if _, err := file.Write(chunk.Data); err != nil {
		log.GetLogger().Fatal("Failed to write chunk to file: ", err.Error())
	}
	if err := file.Close(); err != nil {
		log.GetLogger().Error("Failed to close file while writing chunk: ", err.Error())
	}
}
