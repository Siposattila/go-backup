package backup

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/Siposattila/gobkup/log"
	"github.com/klauspost/compress/zip"
)

type compression struct {
	BackupPath string
	Paths      *[]string
	Exclude    *[]string
}

type filterFunc func(fs.DirEntry) bool

func (c *compression) zipCompress(name string) error {
	var zipFile, err = os.Create(path.Join(c.BackupPath, name))
	if err != nil {
		log.GetLogger().Fatal(err.Error())
	}

	var writer = zip.NewWriter(zipFile)
	for _, path := range *c.Paths {
		if files, err := c.getFiles(path); err == nil {
			c.writeFiles(path, files, writer)
		} else {
			continue
		}
	}

	return writer.Close()
}

// TODO: restore functionality
// func (c *compression) zipDecompress(name string) error {
// 	helper := strings.Split(name, "/")
// 	directoryName := strings.Split(helper[len(helper)-1], ".")[0]
// 	if err := os.Mkdir(directoryName, os.ModePerm); err != nil {
// 		log.GetLogger().Fatal(err.Error())
// 	}

// 	var reader, err = zip.OpenReader(name)
// 	if err != nil {
// 		log.GetLogger().Error(err.Error())
// 	}
// 	c.readFiles(directoryName, reader)

// 	return reader.Close()
// }

// func (c *compression) readFiles(path string, reader *zip.ReadCloser) {
// 	for _, file := range reader.File {
// 		log.GetLogger().Debug(file.Name)
// 		// openedFile, err := file.Open()
// 		// if err != nil {
// 		// 	log.GetLogger().Error(err.Error())
// 		// }

// 		// if _, err := io.
// 	}
// }

func (c *compression) writeFiles(fullPath string, files []fs.DirEntry, writer *zip.Writer) {
	for _, file := range files {
		var fileWriter, err = writer.Create(file.Name())
		if err != nil {
			log.GetLogger().Fatal(err.Error())
		}

		if file.IsDir() {
			dir, _ := c.getFiles(path.Join(fullPath, file.Name()))
			c.writeFiles(fullPath, dir, writer)
		}

		if !file.IsDir() {
			if strings.Contains(fullPath, file.Name()) {
				helper := strings.Split(fullPath, "/")
				helper = helper[:len(helper)-1]
				fullPath = strings.Join(helper, "/")
			}

			var openedFile, openError = os.Open(path.Join(fullPath, file.Name()))
			if openError != nil {
				log.GetLogger().Fatal(openError.Error())
			}

			if _, err := io.Copy(fileWriter, openedFile); err != nil {
				log.GetLogger().Fatal(err.Error())
			}
		}
	}
}

func (c *compression) getFiles(path string) ([]fs.DirEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		log.GetLogger().Error(fmt.Sprintf("Path %s was not found!", path))
		return nil, errors.New("getFiles path not found")
	}

	var files []fs.DirEntry
	if !info.IsDir() {
		files = make([]fs.DirEntry, 0)
		files = append(files, fileInfoDirEntry{info})
	} else {
		files, _ = os.ReadDir(path)
	}

	return c.filter(files, func(file fs.DirEntry) bool {
		for _, exclude := range *c.Exclude {
			match, _ := regexp.MatchString(exclude, file.Name())
			if match {
				return false
			}
		}

		return true
	}), nil
}

func (c *compression) filter(files []fs.DirEntry, f filterFunc) (filtered []fs.DirEntry) {
	for _, file := range files {
		if f(file) {
			filtered = append(filtered, file)
		}
	}

	return
}
