package static

import (
	"embed"
	"io/fs"
)

//go:embed *.css *.js
var assets embed.FS

func Assets() (fs.FS, error) {
	return fs.Sub(assets, "static")
}
