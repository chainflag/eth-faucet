package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed public
var static embed.FS

func Dist() http.FileSystem {
	fsys, err := fs.Sub(static, "public")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
