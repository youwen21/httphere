package conf

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var (
	port    = flag.String("port", "", "Port to run the server. 0 for a random port.")
	backend = flag.String("backend", "", "backend server URL.")
	root    = flag.String("root", "", "specify the root.")

	defaultPort    = "8080"
	defaultBackend = "http://127.0.0.1:80"
)

func init() {
	_ = godotenv.Load()
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
