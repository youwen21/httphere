package conf

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type fastCGIConf struct {
	FastCGI        string
	FastCGIProto   string
	FastCGIAddress string
	FastCGIRoot    string
}

var (
	port    = flag.String("port", "", "Port to run the server. 0 for a random port.")
	root    = flag.String("root", "", "specify the root.")
	backend = flag.String("backend", "", "backend server URL.")

	fastCGI fastCGIConf

	defaultPort    = "8080"
	defaultBackend = "http://127.0.0.1:80"
)

// system vars
var (
	fcgi_proto *string = flag.String("fastcgi_proto", "", "tcp or unix")
	fcgi_addr  *string = flag.String("fastcgi_addr", "", "ip:port or unix file")
)

// user vars
var (
	fcgi_root *string = flag.String("fastcgi_root", "", "the static files root directory, (default empty)")
)

func init() {
	_ = godotenv.Load()

	fastCGI.FastCGI = os.Getenv("FAST_CGI")

	if *fcgi_proto == "" {
		fastCGI.FastCGIProto = os.Getenv("FAST_CGI_PROTO")
	} else {
		fastCGI.FastCGIProto = *fcgi_proto
	}

	if *fcgi_addr == "" {
		fastCGI.FastCGIAddress = os.Getenv("FAST_CGI_ADDR")
	} else {
		fastCGI.FastCGIAddress = *fcgi_addr
	}

	if *fcgi_root == "" {
		fastCGI.FastCGIRoot = os.Getenv("FAST_CGI_ROOT")
	} else {
		fastCGI.FastCGIRoot = *fcgi_root
	}

	if fastCGI.FastCGIRoot == "" {
		fastCGI.FastCGIRoot, _ = os.Getwd()
	}
}

func GetRoot() string {
	if *root != "" {
		return *root
	}

	if rootHere := os.Getenv("ROOT"); rootHere != "" {
		return rootHere
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Getwd err : %v", err)
		os.Exit(1)
	}
	return pwd
}

func GetPort() string {
	if *port != "" {
		return *port
	}

	if ePort := os.Getenv("PORT"); ePort != "" {
		return ePort
	}

	return defaultPort
}

func GetBackend() string {
	if *backend != "" {
		return *backend
	}

	if eBackend := os.Getenv("BACKEND"); eBackend != "" {
		return eBackend
	}

	return defaultBackend
}

func GetFastCGIConfig() fastCGIConf {
	return fastCGI
}
