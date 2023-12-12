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
