package params

const (
	ADDED   int = 1
	REMOVED int = 0
	REMAIN  int = 2
)

type ModifyFileParams struct {
	Path    string   `json:"path"`
	Addr    string   `json:"addr"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Count     int `json:"count"`
	Operation int `json:"operation"`
}
