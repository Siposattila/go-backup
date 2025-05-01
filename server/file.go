package server

import (
	"errors"
	"fmt"
	"os"

	"github.com/Siposattila/gobkup/client"
	"github.com/Siposattila/gobkup/log"
)

const TEMP_FILE = "%s.temp"

func (s *server) writeChunk(chunk *client.Chunk) {
	var file *os.File
	defer file.Close()

	name := fmt.Sprintf(TEMP_FILE, chunk.Name)
	if _, err := os.Stat(fmt.Sprintf("%s/%s", s.Config.BackupPath, name)); errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(name)
		if err != nil {
			log.GetLogger().Error(err)
		}
	} else {
		file, _ = os.Open(name)
	}

	file.Write(chunk.Data)
}
