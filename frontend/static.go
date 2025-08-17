package frontend

import (
	"embed"
	"io/fs"
)

//go:embed dist
var dist embed.FS

var Static, _ = fs.Sub(dist, "dist")
