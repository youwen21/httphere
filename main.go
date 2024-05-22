//
// httphere
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: youwen (youwen21@yeah.net)
//

// https://www.integralist.co.uk/posts/golang-reverse-proxy/

package main

import (
	"flag"
	"fmt"
	"httphere/conf"
	"httphere/server"
	"httphere/utils"
	"net"
	"net/http"
	"os"
)

var (
	methods = `POST, OPTIONS, GET, PUT, DELETE`
	headers = `Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,Accept,Origin,Men,Cache-Control,X-Requested-With,Name,DNT,HOST,Pragma,Referer,Duo,Range,user-Agent,token`
)

func RawCors(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("X-ORIGIN")
		if origin == "" {
			origin = r.Header.Get("ORIGIN")
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin")
		w.Header().Set("Access-Control-Allow-Methods", methods)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", headers)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func printHost() {
	fmt.Println(fmt.Sprintf("%v%v:%v", "http://", "locahost", conf.GetPort()))

	viewUrl, uploadUrl, qrUrl := utils.GetUrls()
	fmt.Println("view url:", viewUrl)
	fmt.Println("upload url:", uploadUrl)
	fmt.Println("qr url:", qrUrl)
}

func main() {
	flag.Parse()

	printHost()

	httpServer := server.NewMyServer()

	addr := net.JoinHostPort(conf.GetHost(), conf.GetPort())
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Listen err: %v", err)
		os.Exit(2)
	}

	if conf.Here.Tls.CertFile != "" && conf.Here.Tls.KeyFile != "" {
		err = http.ServeTLS(listener, RawCors(httpServer), conf.Here.Tls.CertFile, conf.Here.Tls.KeyFile)
		if err != nil {
			fmt.Printf("server http error:%v\n", err)
		}
	} else {
		err = http.Serve(listener, RawCors(httpServer))
		if err != nil {
			fmt.Printf("server http error:%v\n", err)
		}
	}
}
