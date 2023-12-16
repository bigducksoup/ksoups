package message

import "time"

type DirRead struct {
	Path     string `json:"path"`
	FileOnly bool   `json:"fileOnly"`
}

type DirResponse struct {
	Parent   string    `json:"parent"`
	FileOnly bool      `json:"fileOnly"`
	Items    []DirItem `json:"items"`
}

type DirItem struct {
	Name    string    `json:"name"`
	IsDir   bool      `json:"isDir"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"modTime"`
}

type FileRead struct {
	Path string `json:"path"`
}

type FileReadResponse struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int64  `json:"size"`
}

type FileModify struct {
	Path    string   `json:"path"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Count     int `json:"count"`
	Operation int `json:"operation"`
}

const (
	ADDED   int = 1
	REMOVED int = 0
	REMAIN  int = 2
)

type FileModifyResponse struct {
	Path string `json:"path"`
	OK   bool   `json:"ok"`
}
