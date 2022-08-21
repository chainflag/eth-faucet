package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var static embed.FS

func Dist() http.FileSystem {
	fsys, err := fs.Sub(static, "dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
