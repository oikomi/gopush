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
	"encoding/json"
	"github.com/funny/link"
	"github.com/oikomi/gopush/session_manager/redis_store"
)

var InputConfFile = flag.String("conf_file", "session_manager.json", "input conf file name")   

func main() {
	flag.Parse()
	cfg, err := LoadConfig(*InputConfFile)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	
	protocol := link.PacketN(2, binary.BigEndian)
	
	server, err := link.Listen(cfg.TransportProtocols, cfg.Listen, protocol)
	if err != nil {
		panic(err)
	}
	log.Println("server start:", server.Listener().Addr().String())
	
	redisOptions := redis_store.RedisStoreOptions {
			Network :   "tcp",
			Address :   ":6379",
			Database :  1,
			KeyPrefix : "sess",
	}

	redisStore := redis_store.NewRedisStore(&redisOptions)
	
	server.AcceptLoop(func(session *link.Session) {
	log.Println("client", session.Conn().RemoteAddr().String(), "in")

	session.ReadLoop(func(msg link.InMessage) {
		log.Println("client", session.Conn().RemoteAddr().String(),"say:", string(msg))
		
		var ss redis_store.StoreSession
		
		err := json.Unmarshal(msg, &ss)
		if err != nil {
			log.Fatalln("error:", err)
		}

		redisStore.Set(&ss)
	})

	log.Println("client", session.Conn().RemoteAddr().String(), "close")
	})
}
