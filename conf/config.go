package conf

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type FastCGIConf struct {
	FastCGI        string
	FastCGIProto   string
	FastCGIAddress string
	FastCGIRoot    string
}

var (
	defaultPort    = "8090"
	defaultBackend = "http://127.0.0.1:80"

	port    = flag.String("port", "", "Port to run the server. 0 for a random port.")
	root    = flag.String("root", "", "specify the root.")
	backend = flag.String("backend", "", "backend server URL.")

	fastCGI FastCGIConf

	rewrite map[string]string
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

	initFast()
	initRewrite()
}

func initFast() {
	fastCGI.FastCGI = os.Getenv("FASTCGI")

	if *fcgi_proto == "" {
		fastCGI.FastCGIProto = os.Getenv("FASTCGI_PROTO")
	} else {
		fastCGI.FastCGIProto = *fcgi_proto
	}

	if *fcgi_addr == "" {
		fastCGI.FastCGIAddress = os.Getenv("FASTCGI_ADDR")
	} else {
		fastCGI.FastCGIAddress = *fcgi_addr
	}

	if *fcgi_root == "" {
		fastCGI.FastCGIRoot = os.Getenv("FASTCGI_ROOT")
	} else {
		fastCGI.FastCGIRoot = *fcgi_root
	}

	if fastCGI.FastCGIRoot == "" {
		fastCGI.FastCGIRoot, _ = os.Getwd()
	}
}

func initRewrite() {
	rewriteRules := os.Getenv("REWRITE")
	if rewriteRules == "" {
		return
	}
	rewrite = make(map[string]string)

	rules := strings.Split(rewriteRules, ",")
	for _, v := range rules {
		if !strings.Contains(v, "=") {
			fmt.Println("rewrite rule invalid: " + v)
			continue
		}

		rule := strings.Split(v, "=")
		rewrite[rule[0]] = rule[1]
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

func GetFastCGIConfig() FastCGIConf {
	return fastCGI
}

func GetRewrite() map[string]string {
	return rewrite
}
