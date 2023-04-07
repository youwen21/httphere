package server

import (
	"fmt"
	"github.com/youwen21/httphere/conf"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MyServer struct {
	root    string
	rewrite map[string]string

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
		if f.rewrite != nil {
			r.URL.Path = f.RewritePath(r.URL.Path)
		}
		if f.reverseServer != nil {
			f.reverseServer.ServeHTTP(w, r)
		} else {
			b, _ := httputil.DumpRequest(r, true)
			w.Write(b)
			fmt.Println(time.Now().Format(time.RFC3339))
			fmt.Println(string(b))
			fmt.Println()
		}
		return
	}

	// 输出静态文件
	//r.URL.Path = upath
	f.fileServer.ServeHTTP(w, r)
}

func (f MyServer) RewritePath(path string) string {
	if f.rewrite == nil {
		return path
	}

	for k, v := range f.rewrite {
		if strings.HasPrefix(path, k) {
			return strings.Replace(path, k, v, 1)
		}
	}

	return path
}

func NewMyServer() MyServer {
	root := conf.GetRoot()
	backend := conf.GetBackend()

	fmt.Printf("root is %s\n", root)
	fmt.Printf("backend URL is %s\n", backend)

	var s MyServer
	s.root = root
	s.rewrite = conf.GetRewrite()

	s.fileServer = http.FileServer(http.Dir(root))

	if backend != "" {
		backendURL, err := url.Parse(backend)
		if err != nil {
			fmt.Printf("backend server invalid: %v\n", err)
		} else {
			s.reverseServer = httputil.NewSingleHostReverseProxy(backendURL)
		}
	}

	cfg := conf.GetFastCGIConfig()
	if cfg.FastCGI == "true" {
		s.FastCGIEnabled = true
		s.fastCGIServer = NewFastCGIServer(cfg.FastCGIProto, cfg.FastCGIAddress, cfg.FastCGIRoot)
	}

	return s
}
