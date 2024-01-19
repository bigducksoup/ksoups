package param

type CreateShortcutParams struct {
	ProbeId     string `json:"probeId" binding:"required" msg:"probeId 不能为空"`
	Name        string `json:"name" binding:"required" msg:"Name 不能为空"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Timeout     int64  `json:"timeout" binding:"required" msg:"Timeout 不能为空"`
	JustRun     bool   `json:"justRun"`
	Payload     string `json:"payload" binding:"required" msg:"Payload 不能为空"`
}
