package server

import (
	"errors"
	"fmt"
	"os"

	"github.com/Siposattila/go-backup/proto"
)

const TEMP_FILE = "%s.temp"

func (s *server) writeChunk(chunk *proto.BackupChunk) error {
	if _, err := os.Stat(s.Config.BackupPath); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(s.Config.BackupPath, 0777); err != nil {
			return fmt.Errorf("falied to mkdir server backup path: %s", err.Error())
		}
	}

	name := fmt.Sprintf("%s/%s", s.Config.BackupPath, fmt.Sprintf(TEMP_FILE, chunk.Name))
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %s\n%s", name, err.Error())
	}

	if _, err := file.Write(chunk.Data); err != nil {
		return fmt.Errorf("failed to write chunk to file: %s\n%s", name, err.Error())
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close file while writing chunk: %s\n%s", name, err.Error())
	}

	return nil
}
