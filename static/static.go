package static

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed dist/favicon.ico dist/index.html
var RootFS embed.FS

//go:embed dist/assets
var AssetsFS embed.FS

func GetFS(f embed.FS, prefix string) http.FileSystem {
	nf, err := fs.Sub(f, fmt.Sprintf("dist/%s", prefix))
	if err != nil {
		panic(err.Error())
	}
	return http.FS(nf)
}
