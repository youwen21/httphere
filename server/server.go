package server

import (
	"fmt"
	"httphere/conf"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MyServer struct {
	root string

	fileServer       http.Handler
	HistoryRouters   conf.VueHistoryRouters
	domainRevServers map[string]http.Handler
}

func (f MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// domain reverse servers
	if v, has := f.domainRevServers[r.Host]; has {
		v.ServeHTTP(w, f.RewriteRequest(r, conf.Here.GetHostRewrite(r.Host)))
		return
	}

	if r.URL.Path == "/httphere_upload" {
		ServeUpload(w, r)
		return
	}

	if r.URL.Path == "/qr" {
		ServeQr(w, r)
		return
	}

	if f.fileServer != nil {
		// 输出静态文件
		if _, err := os.Stat(filepath.Join(f.root, r.URL.Path)); err == nil {
			f.fileServer.ServeHTTP(w, r)
			return
		}

		if f.HistoryRouters.IsContain(r.URL.Path) {
			fi, _ := os.Open(filepath.Join(f.root, "index.html"))
			content, _ := io.ReadAll(fi)
			w.Write(content)
			return
		}
	}

	// 转发到默认代理服务器
	if v, has := f.domainRevServers["default"]; has {
		v.ServeHTTP(w, f.RewriteRequest(r, conf.Here.GetHostRewrite("default")))
		return
	}

	if conf.Here.Base.DumpRequest == "yes" {
		// dump request
		b, _ := httputil.DumpRequest(r, true)
		fmt.Println(time.Now().Format(time.RFC3339))
		fmt.Println(string(b))
		fmt.Println()
		w.Write(b)
	}
	return
}

func (f MyServer) RewriteRequest(r *http.Request, reMap map[string]string) *http.Request {
	if reMap == nil {
		return r
	}

	for k, v := range reMap {
		if strings.HasPrefix(r.URL.Path, k) {
			r.URL.Path = strings.Replace(r.URL.Path, k, v, 1)
			return r
		}
	}

	return r
}

func initServerByConf(hostConf conf.HostCfg) http.Handler {
	server := http.NewServeMux()
	for cfPath, cfTarget := range hostConf.Paths {
		backendURL, err := url.Parse(cfTarget)
		if err != nil {
			fmt.Println("config parse backend error:", hostConf, cfTarget)
			continue
		}
		if hostConf.ReverseType == "fake_host" {
			server.Handle(cfPath, NewSingleHostReverseProxyFake(backendURL))
		} else {
			server.Handle(cfPath, httputil.NewSingleHostReverseProxy(backendURL))
		}
	}

	return server
}

func NewMyServer() MyServer {
	var s MyServer
	s.domainRevServers = make(map[string]http.Handler)

	// static server
	if conf.Here.Base.StaticServer == "open" {
		root := conf.GetRoot()
		absRoot, _ := filepath.Abs(root)
		fmt.Printf("root is %s\n", absRoot)
		s.root = root
		s.fileServer = http.FileServer(http.Dir(root))
		s.HistoryRouters = conf.Here.Base.HistoryRouters
	}

	for _, v := range conf.Here.Hosts {
		s.domainRevServers[v.Host] = initServerByConf(v)
	}

	return s
}
