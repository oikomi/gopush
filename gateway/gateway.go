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
	"github.com/funny/link"
	"math/rand"
	"time"
	"strconv"
)

var InputConfFile = flag.String("conf_file", "gateway.json", "input conf file name")   

type SessionStore struct {
	ClientID string
	ClientAddr string
	MsgServerAddr string
	ID string
	MaxAge time.Duration
}

func NewSessionStore() *SessionStore {
	return &SessionStore{}
}

func (self *SessionStore)checkClientID(clientID string) bool {
	
	
	return true
}

func selectServer(serverList []string, serverNum int) string{
	return serverList[rand.Intn(serverNum)]
}

func connectSessionManagerServer(cfg Config) (*link.Session, error) {
	protocol := link.PacketN(2, binary.BigEndian)
	client, err := link.Dial("tcp", selectServer(cfg.SessionManagerServerList, len(cfg.SessionManagerServerList)), protocol)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	return client, err
}

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
	sessionManager, err := connectSessionManagerServer(cfg)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	
	defer sessionManager.Close(nil)

	server.AcceptLoop(func(session *link.Session) {
		log.Println("client", session.Conn().RemoteAddr().String(), "in")
		msgServer := selectServer(cfg.MsgServerList, cfg.MsgServerNum)
		
		inMsg, err := session.Read()
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println(string(inMsg))
		
		err = session.Send(link.Binary(msgServer))
		if err != nil {
			log.Fatal(err.Error())
		}

		sessionStore := NewSessionStore()
		sessionStore.ClientID = string(inMsg)
		sessionStore.ClientAddr = session.Conn().RemoteAddr().String()
		sessionStore.MsgServerAddr = msgServer
		sessionStore.ID = strconv.FormatUint(session.Id(), 10)
		
		err = sessionManager.Send(link.JSON {
			sessionStore,
			0,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
		session.Close(nil)
		log.Println("client", session.Conn().RemoteAddr().String(), "close")
	})
}
