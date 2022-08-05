package bplustree

import (
	"os"
)

const (
	PageSize = 4096
)

// DiskManager...
type DiskManager struct {
	file       *os.File
	nextPageID int64
	pageSize   int64
}

// NewNewDiskManager ...
func NewDiskManager(path string) (*DiskManager, error) {
	f, err := CreateOrOpenFile(path)
	if err != nil {
		return nil, err
	}

	metadata, err := GetFileMetadata(f)
	if err != nil {
		return nil, err
	}
	nextPageID := metadata.Size() / PageSize

	return &DiskManager{
		file:       f,
		nextPageID: nextPageID,
		pageSize:   PageSize,
	}, nil
}

// CreateOrOpenFile ...
func CreateOrOpenFile(path string) (*os.File, error) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		f, err := os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// GeGetFileMetadata ...
func GetFileMetadata(file *os.File) (os.FileInfo, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return info, nil
}

// Read ...
func (d *DiskManager) Read(pageID int64) ([]byte, error) {
	offset := d.pageSize * pageID
	// 0 means relative to the origin of the file
	_, err := d.file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, d.pageSize)
	_, err = d.file.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

// Write ...
func (d *DiskManager) Write(pageID int64, data []byte) error {
	offset := d.pageSize * pageID
	// 0 means relative to the origin of the file
	_, err := d.file.Seek(offset, 0)
	if err != nil {
		return err
	}

	_, err = d.file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// Close ...
func (d *DiskManager) Close() error {
	err := d.file.Close()
	if err != nil {
		return err
	}
	return nil
}
