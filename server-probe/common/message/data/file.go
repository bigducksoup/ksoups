package data

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
	Count     int      `json:"count"`
	Operation int      `json:"operation"`
	Value     []string `json:"value"`
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

type FileCreate struct {
	Path       string `json:"path"`
	Permission string `json:"permission"`
}

type FileCreateResponse struct {
	Ok         bool   `json:"ok"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
}

type FileDelete struct {
	Path string `json:"path"`
}

type FileDeleteResponse struct {
	Ok   bool   `json:"ok"`
	Path string `json:"path"`
}
