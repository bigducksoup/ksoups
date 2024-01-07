package data

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
	Name       string    `json:"name"`
	IsDir      bool      `json:"isDir"`
	IsLink     bool      `json:"isLink"`
	LinkTo     string    `json:"linkTo"`
	Size       int64     `json:"size"`
	Permission string    `json:"permission"`
	User       string    `json:"user"`
	UserGroup  string    `json:"usergroup"`
	Mode       string    `json:"mode"`
	ModTime    time.Time `json:"modTime"`
}

type DirCreate struct {
	Path       string `json:"path"`
	Permission string `json:"permission"`
}

type DirCreateResponse struct {
	Ok         bool   `json:"ok"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
}
