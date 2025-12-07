package server

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/Siposattila/go-backup/backup"
	"github.com/Siposattila/go-backup/proto"
)

const (
	MAX_ALLOWED_BACKUP_PER_CLIENT = 5
	TEMP_FILE                     = "%s.temp"
)

func (s *server) writeChunk(chunk *proto.BackupChunk, clientId string) error {
	if _, err := os.Stat(s.Config.BackupPath); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(s.Config.BackupPath, 0644); err != nil {
			return fmt.Errorf("falied to mkdir server backup path: %s", err.Error())
		}
	}

	if _, err := os.Stat(path.Join(s.Config.BackupPath, clientId)); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(path.Join(s.Config.BackupPath, clientId), 0644); err != nil {
			return fmt.Errorf("falied to mkdir server backup path for client: %s", err.Error())
		}
	}

	name := path.Join(s.Config.BackupPath, clientId, fmt.Sprintf(TEMP_FILE, chunk.Name))
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

// TODO: this is temporary should implement a more configurable system for it
func (s *server) applyRetentionPolicy(clientId string) error {
	files, readDirError := os.ReadDir(path.Join(s.Config.BackupPath, clientId))
	if readDirError != nil {
		return fmt.Errorf("failed to read backup dir for %s", clientId)
	}
	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()

		return infoI.ModTime().Before(infoJ.ModTime())
	})

	amountOfBackups := 0
	for _, file := range files {
		if strings.Contains(file.Name(), backup.BACKUP_EXTENSION) && !strings.Contains(file.Name(), ".temp") {
			amountOfBackups++
		}
	}

	if amountOfBackups >= MAX_ALLOWED_BACKUP_PER_CLIENT {
		if removeError := os.Remove(path.Join(s.Config.BackupPath, clientId, files[0].Name())); removeError != nil {
			return fmt.Errorf("failed to carry out retention policy on %s old backup %s", clientId, files[0].Name())
		}
	}

	return nil
}
