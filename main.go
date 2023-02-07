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

package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"personal/httphere/conf"
	"personal/httphere/server"
)

func main() {
	sRoot := conf.GetRoot()
	sPort := conf.GetPort()
	sBackend := conf.GetBackend()

	httpServer := server.NewMyServer(sRoot, sBackend)
	addr := net.JoinHostPort("", sPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Listen err: %v", err)
		os.Exit(2)
	}

	fmt.Printf("port is %s\n", sPort)
	fmt.Printf("backend URL is %s\n", sBackend)
	fmt.Printf("root is %s\n", sRoot)

	fmt.Printf("Listening on %s\n", listener.Addr().String())
	err = http.Serve(listener, httpServer)
	if err != nil {
		fmt.Printf("server http error:%v\n", err)
	}
}
