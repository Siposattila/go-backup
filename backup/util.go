package backup

import "io/fs"

type fileInfoDirEntry struct {
	fs.FileInfo
}

func (f fileInfoDirEntry) Type() fs.FileMode {
	return f.Mode().Type()
}

func (f fileInfoDirEntry) Info() (fs.FileInfo, error) {
	return f.FileInfo, nil
}
