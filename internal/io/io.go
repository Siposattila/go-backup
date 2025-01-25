package io

import (
	"errors"
	"os"

	"github.com/Siposattila/gobkup/internal/console"
)

func CreateDir(path string) {
	if _, statError := os.Stat(path); errors.Is(statError, os.ErrNotExist) {
		var mkdirError = os.Mkdir(path, os.ModePerm)
		if mkdirError != nil {
			console.Fatal("Unable to create dir: " + mkdirError.Error())
		}
	}
}

func WriteFile(path string, fileName string, data []byte) {
	var writeError = os.WriteFile(path+"/"+fileName, data, 0644)
	if writeError != nil {
		console.Fatal("Unable to write to file: " + writeError.Error())
	}
}

func ReadFile(path string, fileName string) []byte {
	var buffer, readError = os.ReadFile(path + "/" + fileName)
	if readError != nil {
		console.Fatal("Unable to read file: " + readError.Error())
	}

	return buffer
}
