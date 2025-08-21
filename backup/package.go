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

	"github.com/Siposattila/go-backup/log"
	"github.com/klauspost/compress/zip"
)

type compression struct {
	BackupPath string
	Paths      *[]string
	Exclude    *[]string
	Store      *store
}

type filterFunc func(fs.DirEntry) bool

func (c *compression) zipCompress(name string) (zipPath string) {
	zipPath = path.Join(c.BackupPath, name)
	var zipFile, err = os.Create(zipPath)
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
	if err := writer.Close(); err != nil {
		log.GetLogger().Fatal(err.Error())
	}

	r, _ := zip.OpenReader(zipPath)
	if len(r.File) == 0 {
		log.GetLogger().Warning("Empty archive!")
		if err := os.Remove(zipPath); err != nil {
			log.GetLogger().Error(err.Error())
		}

		zipPath = ""
	}
	if err := r.Close(); err != nil {
		log.GetLogger().Fatal(err.Error())
	}

	return
}

func (c *compression) writeFiles(fullPath string, files []fs.DirEntry, writer *zip.Writer) {
	for _, file := range files {
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

			checksum := c.Store.checksum(path.Join(fullPath, file.Name()))
			if contains, index := c.Store.contains(file.Name()); (contains &&
				checksum != c.Store.Entries[index].Checksum) || !contains {
				if contains {
					c.Store.Entries[index].Checksum = checksum
				} else {
					c.Store.add(entry{Name: file.Name(), Checksum: checksum})
				}

				if fileWriter, err := writer.Create(file.Name()); err == nil {
					if _, err := io.Copy(fileWriter, openedFile); err != nil {
						log.GetLogger().Fatal(err.Error())
					}
				} else {
					log.GetLogger().Fatal(err.Error())
				}
			}
		}
	}
	c.Store.saveStore()
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
