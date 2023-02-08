package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"personal/httphere/conf"
	"strings"
)

type MyServer struct {
	root string

	fileServer    http.Handler
	reverseServer *httputil.ReverseProxy

	FastCGIEnabled bool
	fastCGIServer  http.Handler
}

func (f MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}

	// 如果以.php结尾
	if f.FastCGIEnabled && strings.HasSuffix(upath, ".php") {
		f.fastCGIServer.ServeHTTP(w, r)
		return
	}
	// 如果静态文件不存在
	_, err := os.Stat(filepath.Join(f.root, upath))
	if os.IsNotExist(err) {
		f.reverseServer.ServeHTTP(w, r)
		return
	}

	// 输出静态文件
	//r.URL.Path = upath
	f.fileServer.ServeHTTP(w, r)
}

func NewMyServer(root string, proxyURL string) MyServer {
	var s MyServer
	s.root = root

	s.fileServer = http.FileServer(http.Dir(root))

	backendURL, err := url.Parse(proxyURL)
	if err != nil {
		fmt.Printf("backend server invalid: %v\n", err)
	}
	s.reverseServer = httputil.NewSingleHostReverseProxy(backendURL)

	cfg := conf.GetFastCGIConfig()
	if cfg.FastCGI == "true" {
		s.FastCGIEnabled = true
		s.fastCGIServer = NewFastCGIServer(cfg.FastCGIProto, cfg.FastCGIAddress, cfg.FastCGIRoot)
	}

	return s
}
