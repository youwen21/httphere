package server

import (
	"fmt"
	"httphere/conf"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type MyServer struct {
	root string

	fileServer       http.Handler
	domainRevServers map[string]http.Handler
}

func (f MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// domain servers
	if v, has := f.domainRevServers[r.Host]; has {
		v.ServeHTTP(w, r)
		return
	}

	if f.fileServer != nil {
		_, err := os.Stat(filepath.Join(f.root, r.URL.Path))
		// 输出静态文件
		if err == nil {
			f.fileServer.ServeHTTP(w, r)
			return
		}
	}

	// 转发到默认代理服务器
	if v, has := f.domainRevServers["default"]; has {
		v.ServeHTTP(w, r)
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

func initMuxServerByConf(hostConf conf.HostConf) *http.ServeMux {
	server := http.NewServeMux()
	for k, v := range hostConf.Paths {
		backendURL, err := url.Parse(v)
		if err != nil {
			fmt.Println("config has error:", hostConf)
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
		s.domainRevServers[v.Host] = initMuxServerByConf(v)
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
