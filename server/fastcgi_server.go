package server

// TODO:  use net/http/fcgi

// copy from https://github.com/alash3al/xerver/blob/master/xerver.go

import (
	"crypto/tls"
	"fmt"
	"httphere/fcgiclient"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type FastCGIServer struct {
	FastCGIProto   string
	FastCGIAddress string
	FastCGIRoot    string
}

func (f FastCGIServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fcgi, err := fcgiclient.Dial(f.FastCGIProto, f.FastCGIAddress)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Unable to connect to the backend", 502)
		return
	}
	// automatically close the fastcgi connection and the requested body at the end .
	defer fcgi.Close()
	defer req.Body.Close()

	reqHost, reqPort, err := net.SplitHostPort(req.Host)
	if err != nil {
		// whatever, just assume there was no port
		reqHost = req.Host
	}

	remoteAddr, remotePort, _ := net.SplitHostPort(req.RemoteAddr)
	req.URL.Path = req.URL.ResolveReference(req.URL).Path

	requestScheme := "http"
	if req.TLS != nil {
		requestScheme = "https"
	}

	// TODO
	authUser := ""

	env := map[string]string{
		// from https://github.com/caddyserver/caddy/blob/master/modules/caddyhttp/reverseproxy/fastcgi/fastcgi.go
		// Variables defined in CGI 1.1 spec
		"AUTH_TYPE":         "",
		"CONTENT_LENGTH":    fmt.Sprintf("%d", req.ContentLength),
		"CONTENT_TYPE":      req.Header.Get("Content-Type"),
		"GATEWAY_INTERFACE": "CGI/1.1",
		"PATH_INFO":         req.URL.Path,
		"QUERY_STRING":      req.URL.Query().Encode(),

		"REMOTE_ADDR":  remoteAddr,
		"REMOTE_HOST":  remoteAddr,
		"REMOTE_PORT":  remotePort,
		"REMOTE_IDENT": "",
		"REMOTE_USER":  authUser,

		"REQUEST_METHOD":  req.Method,
		"REQUEST_SCHEME":  requestScheme,
		"SERVER_NAME":     reqHost,
		"SERVER_PROTOCOL": req.Proto,
		"SERVER_SOFTWARE": "",

		"DOCUMENT_ROOT":   f.FastCGIRoot,
		"DOCUMENT_URI":    req.URL.Path,
		"HTTP_HOST":       req.Host, // added here, since not always part of headers
		"REQUEST_URI":     req.URL.RequestURI(),
		"SCRIPT_FILENAME": filepath.Join(f.FastCGIRoot, req.URL.Path),
		"SCRIPT_NAME":     req.URL.Path,

		// 以下这几个 好像没什么用？
		"REQUEST_PATH":  req.URL.Path,
		"SERVER_ADDR":   reqHost,
		"SERVER_PORT":   reqPort,
		"FCGI_PROTOCOL": f.FastCGIProto,
		"fcgi_addr":     f.FastCGIAddress,
	}

	if reqPort != "" {
		env["SERVER_PORT"] = reqPort
	} else if requestScheme == "http" {
		env["SERVER_PORT"] = "80"
	} else if requestScheme == "https" {
		env["SERVER_PORT"] = "443"
	}

	// tell fastcgi backend that, this connection is done over https connection if enabled .
	if req.TLS != nil {
		env["HTTPS"] = "on"
		v, ok := tlsProtocolStrings[req.TLS.Version]
		if ok {
			env["SSL_PROTOCOL"] = v
		}
		// TODO
		//env["SSL_CIPHER"] = cs.Name
	}
	// iterate over request headers and append them to the environment varibales in the valid format .
	for k, v := range req.Header {
		env["HTTP_"+strings.Replace(strings.ToUpper(k), "-", "_", -1)] = strings.Join(v, ";")
	}
	// fethcing the response from the fastcgi backend,
	// and check for errors .
	resp, err := fcgi.Request(env, req.Body)
	if err != nil {
		fmt.Println("err> ", err.Error())
		http.Error(res, "Unable to fetch the response from the backend", 502)
		return
	}
	// parse the fastcgi status .
	resp.Status = resp.Header.Get("Status")
	resp.StatusCode, _ = strconv.Atoi(strings.Split(resp.Status, " ")[0])
	if resp.StatusCode < 100 {
		resp.StatusCode = 200
	}
	// automatically close the fastcgi response body at the end .
	defer resp.Body.Close()
	// read the fastcgi response headers,
	// exclude "Xerver-Internal-*" headers from the response,
	// and apply the actions related to them .
	for k, v := range resp.Header {
		if !strings.HasPrefix(k, "Xerver-Internal-") {
			for i := 0; i < len(v); i++ {
				if res.Header().Get(k) == "" {
					res.Header().Set(k, v[i])
				} else {
					res.Header().Add(k, v[i])
				}
			}
		}
	}
	// remove server tokens from the response
	if resp.Header.Get("Xerver-Internal-ServerTokens") != "off" {
		res.Header().Set("Server", "httphere")
	}
	// serve the provided filepath using theinternal fileserver
	if resp.Header.Get("Xerver-Internal-FileServer") != "" {
		res.Header().Del("Content-Type")
		http.ServeFile(res, req, resp.Header.Get("Xerver-Internal-FileServer"))
		return
	}
	// serve the response from another backend "http-proxy"
	if resp.Header.Get("Xerver-Internal-ProxyPass") != "" {
		u, e := url.Parse(resp.Header.Get("Xerver-Internal-ProxyPass"))
		if e != nil {
			fmt.Println("err> ", e.Error())
			http.Error(res, "Invalid internal-proxypass value", 502)
			return
		}
		httputil.NewSingleHostReverseProxy(u).ServeHTTP(res, req)
		return
	}
	// fix the redirect issues by fetching the fastcgi response location header
	// then redirect the client, then ignore any output .
	if resp.Header.Get("Location") != "" {
		http.Redirect(res, req, resp.Header.Get("Location"), resp.StatusCode)
		return
	}
	// write the response status code .
	res.WriteHeader(resp.StatusCode)
	// only sent the header if the request method isn't HEAD .
	if req.Method != "HEAD" {
		io.Copy(res, resp.Body)
	}
}

// Map of supported protocols to Apache ssl_mod format
// Note that these are slightly different from SupportedProtocols in caddytls/config.go
var tlsProtocolStrings = map[uint16]string{
	tls.VersionTLS10: "TLSv1",
	tls.VersionTLS11: "TLSv1.1",
	tls.VersionTLS12: "TLSv1.2",
	tls.VersionTLS13: "TLSv1.3",
}

func NewFastCGIServer(proto string, address string, root string) FastCGIServer {
	var s FastCGIServer
	s.FastCGIProto = proto
	s.FastCGIAddress = address
	s.FastCGIRoot = root
	return s
}
