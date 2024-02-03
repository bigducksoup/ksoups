package param

type DirCreateParams struct {
	Path       string `json:"path" binding:"required"`
	ProbeId    string `json:"probeId" binding:"required"`
	Permission string `json:"permission" binding:"required"`
}
