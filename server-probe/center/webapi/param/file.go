package param

const (
	ADDED   int = 1
	REMOVED int = 0
	REMAIN  int = 2
)

type ModifyFileParams struct {
	Path    string   `json:"path"`
	ProbeId string   `json:"probeId"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Count     int      `json:"count"`
	Operation int      `json:"operation"`
	Value     []string `json:"value"`
}

type FileCreateParams struct {
	Path       string `json:"path"`
	ProbeId    string `json:"probeId"`
	Permission string `json:"permission"`
}
