package dto

import (
	"fmt"
	"io/fs"
	"path"
	"time"
)

const (
	FSS3StorageClass = "STANDARD"
)

var DefaultOwner = OwnerInfo{
	ID:          "0",
	DisplayName: "root",
}

type DirInfo struct {
	BucketName  string     `xml:"Name"`
	Path        string     `xml:"Prefix"`
	IsTruncated bool       `xml:"IsTruncated"`
	Files       []FileInfo `xml:"Contents"`
}

type FileInfo struct {
	FileName         string    `xml:"FileName"`
	IsDirectory      bool      `xml:"IsDirectory"`
	FullPath         string    `xml:"Key"`
	LastModified     string    `xml:"LastModified"`
	LastModifiedUnix int64     `xml:"LastModifiedUnix"`
	Hash             string    `xml:"ETag"`
	Size             uint64    `xml:"Size"`
	StorageClass     string    `xml:"StorageClass"`
	Owner            OwnerInfo `xml:"Owner"`
}

type OwnerInfo struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

func NewFSDirInfo(serverName, path string) DirInfo {
	return DirInfo{
		BucketName: serverName,
		Path:       path,
		Files:      []FileInfo{},
	}
}

func (i *DirInfo) AddFSFile(f fs.FileInfo) {
	name := f.Name()
	size := uint64(f.Size())
	isDir := f.IsDir()
	lastModified := f.ModTime()

	if isDir {
		// Linux stat returns directory size, but we don't want it
		size = 0
	}

	file := FileInfo{
		FileName:         name,
		IsDirectory:      f.IsDir(),
		FullPath:         path.Join(i.Path, name),
		LastModified:     lastModified.UTC().Format(time.RFC3339),
		LastModifiedUnix: lastModified.Unix(),
		Hash:             FastHash([]byte(fmt.Sprintf("%s %s %d %d", i.Path, name, size, lastModified.UnixNano()))),
		Size:             size,
		StorageClass:     FSS3StorageClass,
		Owner:            DefaultOwner,
	}
	i.Files = append(i.Files, file)
}

func FastHash(v []byte) string {
	hi := uint64(0x66ccff9920120712)
	lo := uint64(0x114514190d000721)
	for _, b := range v {
		hi ^= uint64(b)
		hi ^= hi << 13
		hi ^= lo >> 7
		hi ^= hi << 17
		lo ^= hi
		lo ^= hi << 19
		lo ^= lo >> 5
		lo ^= hi << 11
	}
	return fmt.Sprintf("%016x%016x", hi, lo)
}
