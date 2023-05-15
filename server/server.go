package server

import (
	"fmt"
	"github.com/spf13/cast"
	"httphere/conf"
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

	fileServer    http.Handler
	reverseServer *httputil.ReverseProxy

	fastCGIServer http.Handler
}

func (f MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	config := conf.GetConfig()
	if v, has := config[r.Host]; has {
		f.ServeHTTPHost(w, r, v)
		return
	}

	upath := r.URL.Path

	_, err := os.Stat(filepath.Join(f.root, upath))
	// 输出静态文件
	if !os.IsNotExist(err) {
		f.fileServer.ServeHTTP(w, r)
		return
	}

	// 转发到默认代理服务器
	if f.reverseServer != nil {
		f.reverseServer.ServeHTTP(w, r)
		return
	}

	// 没有代理服务器，打应请求内容到 日志和请求方
	b, _ := httputil.DumpRequest(r, true)
	fmt.Println(time.Now().Format(time.RFC3339))
	fmt.Println(string(b))
	fmt.Println()
	w.Write(b)
}

func (f MyServer) ServeHTTPHost(w http.ResponseWriter, r *http.Request, hostConf interface{}) {
	upath := r.URL.Path

	rewriteRule := getRewriteRule(hostConf)
	if rewriteRule != nil {
		upath = RewritePath(upath, rewriteRule)
	}

	pathsCfg := getPaths(hostConf)
	backend, has := pathsCfg[upath]
	if !has {
		backend = pathsCfg["backend"]
	}
	backendURL, err := url.Parse(backend)
	if err != nil {
		w.Write([]byte("backend config error"))
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	proxy.ServeHTTP(w, r)
}

func getRewriteRule(cfg interface{}) map[string]string {
	confMap, ok := cfg.(map[string]interface{})
	if !ok {
		return nil
	}
	ruleI, ok := confMap["rewrite"]
	if !ok {
		return nil
	}
	rule, ok := ruleI.(map[string]interface{})
	if !ok {
		return nil
	}

	result := make(map[string]string)
	for k, v := range rule {
		result[k] = cast.ToString(v)
	}
	return result
}

func getPaths(cfg interface{}) map[string]string {
	confMap, ok := cfg.(map[string]interface{})
	if !ok {
		return nil
	}

	result := make(map[string]string)
	for k, v := range confMap {
		if k == "rewrite" {
			continue
		}
		if vStr, vok := v.(string); vok {
			result[k] = vStr
		}
	}

	return result
}

func RewritePath(path string, rewriteRule map[string]string) string {
	for k, v := range rewriteRule {
		if strings.HasPrefix(path, k) {
			return strings.Replace(path, k, v, 1)
		}
	}

	return path
}

//func (f MyServer) RewritePath(path string) string {
//	if f.rewrite == nil {
//		return path
//	}
//
//	for k, v := range f.rewrite {
//		if strings.HasPrefix(path, k) {
//			return strings.Replace(path, k, v, 1)
//		}
//	}
//
//	return path
//}

func NewMyServer() MyServer {
	var s MyServer

	// static server
	root := conf.GetRoot()
	fmt.Printf("root is %s\n", root)
	s.root = root
	s.fileServer = http.FileServer(http.Dir(root))

	// default reverse proxy server
	backend := conf.GetBackend()
	fmt.Printf("backend URL is %s\n", backend)
	if backend != "" {
		backendURL, err := url.Parse(backend)
		if err != nil {
			fmt.Printf("backend server invalid: %v\n", err)
		} else {
			s.reverseServer = httputil.NewSingleHostReverseProxy(backendURL)
		}
	}

	//http.HandleFunc()
	// mux server， support config.toml

	return s
}
