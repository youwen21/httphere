package main

import (
	"net/http"
)

type FastCGIServer struct {
}

// TODO
func (f FastCGIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func NewFastCGIServer() FastCGIServer {
	var s FastCGIServer
	return s
}
