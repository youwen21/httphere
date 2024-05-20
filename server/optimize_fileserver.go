package server

import (
	"httphere/asset"
	"io"
	"net/http"
	"strings"
)

type OptFileSer struct {
	dirRoot string
}

// 这个作毛线用的？
func OptFileServer(next http.Handler, root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpFs := http.Dir(root)

		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}

		f, err := httpFs.Open(upath)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer f.Close()
		d, err := f.Stat()
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if d.IsDir() {
			fi, _ := asset.Dist.Open("dist/view_files_part1.html")
			content, _ := io.ReadAll(fi)
			w.Write(content)
		}

		next.ServeHTTP(w, r)

		if d.IsDir() {
			fi, _ := asset.Dist.Open("dist/view_files_part2.html")
			content, _ := io.ReadAll(fi)
			w.Write(content)
		}
	})
}
