package compression

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/Siposattila/gobkup/internal/console"
	"github.com/klauspost/compress/zip"
	"github.com/klauspost/compress/zstd"
)

type Compression struct {
	Path string
}

func compress(in io.Reader, out io.Writer) error {
	var enc, err = zstd.NewWriter(out)
	if err != nil {
		return err
	}

	_, err = io.Copy(enc, in)
	if err != nil {
		enc.Close()
		return err
	}

	return enc.Close()
}

func (c Compression) writeFiles(files []fs.DirEntry, writer *zip.Writer) {
	for _, file := range files {
		var fileWriter, err = writer.Create(file.Name())
		if err != nil {
			console.Fatal(err.Error())
		}

		if file.IsDir() {
			c.writeFiles(getFiles(path.Join(c.Path, file.Name())), writer)
		}

		if !file.IsDir() {
			var openedFile, openError = os.Open(path.Join(c.Path, file.Name()))
			if openError != nil {
				console.Fatal(openError.Error())
			}
			var fileReader io.Reader
			fileReader = openedFile

			var fileCompressed = new(bytes.Buffer)
			compress(fileReader, fileCompressed)

			_, err = fileWriter.Write(fileCompressed.Bytes())
			if err != nil {
				console.Fatal(err.Error())
			}
		}
	}

	return
}

func getFiles(name string) []fs.DirEntry {
	var files, error = os.ReadDir(name)
	if error != nil {
		console.Fatal("The provided path is not found!")
	}

	return files
}

func (c Compression) ZipCompress(name string) error {
	var zipFile, createError = os.Create(name)
	if createError != nil {
		console.Fatal(createError.Error())
	}

	var writer = zip.NewWriter(zipFile)
	c.writeFiles(getFiles(c.Path), writer)

	return writer.Close()
}

func (c Compression) ZipDecompress(name string) error {
	// TODO: implement
	return errors.New("implement")
}
