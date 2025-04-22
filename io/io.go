package io

import (
	"errors"
	"os"
	"path"

	"github.com/Siposattila/gobkup/log"
)

func CreateDir(path string) {
	if _, statError := os.Stat(path); errors.Is(statError, os.ErrNotExist) {
		mkdirError := os.Mkdir(path, os.ModePerm)
		if mkdirError != nil {
			log.GetLogger().Fatal("Unable to create dir.", mkdirError.Error())
		}
	}
}

func WriteFile(dir, name string, data []byte) {
	writeError := os.WriteFile(path.Join(dir, name), data, 0644)
	if writeError != nil {
		log.GetLogger().Fatal("Unable to write to file.", writeError.Error())
	}
}

func ReadFile(dir, name string) ([]byte, error) {
	buffer, readError := os.ReadFile(path.Join(dir, name))
	if readError != nil {
		return *new([]byte), readError
	}

	return buffer, nil
}

func Delete(name string) {
	if err := os.Remove(name); err != nil {
		log.GetLogger().Warning("Was not able to delete " + name)
	}
}
