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
	"time"
	"encoding/json"
	"github.com/funny/link"
	"github.com/oikomi/gopush/session_manager/redis_store"
	"github.com/oikomi/gopush/protocol"
)

var InputConfFile = flag.String("conf_file", "session_manager.json", "input conf file name")

type SessionManager struct {
	
}   

func connectMsgServer(ms string) (*link.Session, error) {
	p := link.PacketN(2, link.BigEndianBO, link.LittleEndianBF)
	client, err := link.Dial("tcp", ms, p)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	return client, err
}

func handleMsgServerClient(msc *link.Session, redisStore *redis_store.RedisStore) {
	msc.ReadLoop(func(msg link.InBuffer) {
		log.Println("client", msc.Conn().RemoteAddr().String(),"say:", string(msg.Get()))
		
		var ss redis_store.StoreSession
		
		log.Println(string(msg.Get()))
		
		err := json.Unmarshal(msg.Get(), &ss)
		if err != nil {
			log.Fatalln("error:", err)
		}

		err = redisStore.Set(&ss)
		if err != nil {
			log.Fatalln("error:", err)
		}
	})

	log.Println("client", msc.Conn().RemoteAddr().String(), "close")
}

func subscribeChannels(cfg Config, redisStore *redis_store.RedisStore) {
	var msgServerClientList []*link.Session
	for _, ms := range cfg.MsgServerList {
		msgServerClient, err := connectMsgServer(ms)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		cmd := protocol.NewCmd()
		
		cmd.Cmd = protocol.SUBSCRIBE_CHANNEL_CMD
		cmd.Args[0] = SYSCTRL_CLIENT_STATUS
		
		msgServerClient.Send(link.JSON {
			cmd,
		})
		
		msgServerClientList = append(msgServerClientList, msgServerClient)
	}

	for _, msc := range msgServerClientList {
		go handleMsgServerClient(msc, redisStore)
	}
}

func main() {
	flag.Parse()
	cfg, err := LoadConfig(*InputConfFile)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	
	p := link.PacketN(2, link.BigEndianBO, link.LittleEndianBF)
	
	server, err := link.Listen(cfg.TransportProtocols, cfg.Listen, p)
	if err != nil {
		panic(err)
	}
	log.Println("server start:", server.Listener().Addr().String())
	
	redisOptions := redis_store.RedisStoreOptions {
			Network :   "tcp",
			Address :   cfg.Redis.Port,
			ConnectTimeout : time.Duration(cfg.Redis.ConnectTimeout)*time.Millisecond,
			ReadTimeout : time.Duration(cfg.Redis.ReadTimeout)*time.Millisecond,
			WriteTimeout : time.Duration(cfg.Redis.WriteTimeout)*time.Millisecond,
			Database :  1,
			KeyPrefix : "push",
	}

	redisStore := redis_store.NewRedisStore(&redisOptions)

	go subscribeChannels(cfg, redisStore)
	
	server.AcceptLoop(func(session *link.Session) {
	
	})
}
