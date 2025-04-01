package io

import (
	"errors"
	"log"
	"os"
)

func CreateDir(path string) {
	if _, statError := os.Stat(path); errors.Is(statError, os.ErrNotExist) {
		mkdirError := os.Mkdir(path, os.ModePerm)
		if mkdirError != nil {
			log.Fatal("Unable to create dir.", mkdirError.Error())
		}
	}
}

func WriteFile(path, fileName string, data []byte) {
	writeError := os.WriteFile(path+"/"+fileName, data, 0644)
	if writeError != nil {
		log.Fatal("Unable to write to file.", writeError.Error())
	}
}

func ReadFile(path, fileName string) ([]byte, error) {
	buffer, readError := os.ReadFile(path + "/" + fileName)
	if readError != nil {
		return *new([]byte), readError
	}

	return buffer, nil
}
