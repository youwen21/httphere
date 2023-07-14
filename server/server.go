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
	domainRevServers map[string]http.Handler
}

func getHostRewrite(host string) map[string]string {
	for _, v := range conf.Here.Hosts {
		if v.Host == host {
			return v.Rewrite
		}
	}

	return nil
}

func (f MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// domain servers
	if v, has := f.domainRevServers[r.Host]; has {
		v.ServeHTTP(w, f.RewriteRequest(r, getHostRewrite(r.Host)))
		return
	}

	if f.fileServer != nil {
		_, err := os.Stat(filepath.Join(f.root, r.URL.Path))
		// 输出静态文件
		if err == nil {
			f.fileServer.ServeHTTP(w, r)
			return
		}

		if conf.Here.Base.HistoryRouters.IsContain(r.URL.Path) {
			fi, _ := os.Open(filepath.Join(f.root, "index.html"))
			content, _ := io.ReadAll(fi)
			w.Write(content)
			return
		}
	}

	// 转发到默认代理服务器
	if v, has := f.domainRevServers["default"]; has {
		v.ServeHTTP(w, f.RewriteRequest(r, getHostRewrite("default")))
		return
	}

	if conf.Here.Base.DumpRequest == "yes" {
		// 没有代理服务器，打应请求内容到 日志和请求方
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
			path := strings.Replace(r.URL.Path, k, v, 1)
			r.URL.Path = path
			return r
		}
	}

	return r
}

func initServerByConf(hostConf conf.HostConf) http.Handler {
	if hostConf.ReverseType == "fast_cgi" {
		return NewFastCGIServer(hostConf.FastCGI.Proto, hostConf.FastCGI.Address, hostConf.FastCGI.Root)
	}

	server := http.NewServeMux()
	for k, v := range hostConf.Paths {
		backendURL, err := url.Parse(v)
		if err != nil {
			fmt.Println("config parse backend error:", hostConf, v)
			continue
		}
		if hostConf.ReverseType == "fake_host" {
			server.Handle(k, NewSingleHostReverseProxyFake(backendURL))
		} else {
			server.Handle(k, httputil.NewSingleHostReverseProxy(backendURL))
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
		fmt.Printf("root is %s\n", root)
		s.root = root
		s.fileServer = http.FileServer(http.Dir(root))
	}

	for _, v := range conf.Here.Hosts {
		s.domainRevServers[v.Host] = initServerByConf(v)
	}

	// overwrite default reverse server if command backend flag params represent
	backend := conf.GetBackend()
	if backend != "" {
		backendURL, err := url.Parse(backend)
		if err == nil {
			s.domainRevServers["default"] = httputil.NewSingleHostReverseProxy(backendURL)
		} else {
			fmt.Println("backend parse error:", backend, err)
		}
	}

	return s
}
