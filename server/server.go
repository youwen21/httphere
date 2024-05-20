package server

import (
	"bufio"
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
	// domain servers
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
		// 没有代理服务器，打印日志
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
	if hostConf.ReverseType == "fast_cgi" {
		return NewFastCGIServer(hostConf.FastCGI.Proto, hostConf.FastCGI.Address, hostConf.FastCGI.Root)
	}

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
		fmt.Printf("root is %s\n", root)
		s.root = root
		//s.fileServer = OptFileServer(http.FileServer(http.Dir(root)), root)
		s.fileServer = http.FileServer(http.Dir(root))
		// init VueHistoryRouters
		s.HistoryRouters = conf.Here.Base.HistoryRouters
		//if _, err := os.Stat(filepath.Join(root, "httphere_routers.txt")); err == nil {
		//	hsRouters, err := initHistoryRouters(filepath.Join(root, "httphere_routers.txt"))
		//	if err != nil {
		//		fmt.Println("initHistoryRouters err:", err)
		//	} else {
		//		s.HistoryRouters = hsRouters
		//	}
		//}
	}

	for _, v := range conf.Here.Hosts {
		s.domainRevServers[v.Host] = initServerByConf(v)
	}

	return s
}

func initHistoryRouters(routerFile string) (conf.VueHistoryRouters, error) {
	hsRouters := make(conf.VueHistoryRouters, 0)
	file, err := os.Open(routerFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := bufio.NewReader(file)
	for {
		// ReadLine is a low-level line-reading primitive.
		// Most callers should use ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
		bytes, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		hsRouters = append(hsRouters, string(bytes))
	}

	return hsRouters, nil
}
