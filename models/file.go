package models

import (
	"time"
)

const (
	TypeMp3     = "mp3"
	TypeUnknown = "unknown"
)

type Tags struct {
	Title       string
	Artist      string
	Album       string
	PublishedAt string
	Genre       string
}

type File struct {
	Filename string
	AddedAt  time.Time
	Metadata Tags
}

type FileListJson struct {
	Quota   int
	Current int
	List    []File
}
