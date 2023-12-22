package static

import (
	"embed"
	"io/fs"
)

//go:embed *.js *.css
var assets embed.FS

func Assets() (fs.FS, error) {
	return fs.Sub(assets, ".")
}
