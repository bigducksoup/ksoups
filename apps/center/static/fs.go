package static

import "embed"

//go:embed dist/*
//go:embed dist/assets/*
var DistFS embed.FS
