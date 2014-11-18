//
// Copyright 2014 Hong Miao. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"log"
	"encoding/binary"
	"github.com/oikomi/gopush/netlib"
)

var InputConfFile = flag.String("conf_file", "gateway.json", "input conf file name")   

func handler(session *netlib.Session) {
	log.Println("client", session.Conn().RemoteAddr().String(), "in")

	session.ReadLoop(func(msg []byte) {
		log.Println("client", session.Conn().RemoteAddr().String(), "say:", string(msg))
		session.Send(netlib.Binary(msg))
	})

	log.Println("client", session.Conn().RemoteAddr().String(), "close")
	
	
}

func main() {
	flag.Parse()
	//log.Println(*InputConfFile)
	cfg, err := LoadConfig(*InputConfFile)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	
	protocol := netlib.PacketN(2, binary.BigEndian)
	
	server, err := netlib.Listen(cfg.TransportProtocols, cfg.Listen, protocol)
	if err != nil {
		panic(err)
	}
	log.Println("server start:", server.Listener().Addr().String())
	
	server.AcceptLoop(handler)
	//log.Println(server.sessions)

}
