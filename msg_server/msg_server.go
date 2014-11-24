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
	"github.com/funny/link"
	"strconv"
)

var InputConfFile = flag.String("conf_file", "msg_server.json", "input conf file name")   

type MsgServer struct {
	channels ChannelMap
	server *link.Server
}

func NewMsgServer() *MsgServer {
	ms := &MsgServer {
		channels : make(ChannelMap),
	}
	
	return ms
}

func (self *MsgServer)initChannels() {
	channel := link.NewChannel(self.server.Protocol())
	self.channels[SYSCTRL_CLIENT_STATUS] = channel
}

func main() {
	flag.Parse()
	cfg, err := LoadConfig(*InputConfFile)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	
	ms := NewMsgServer()
	
	protocol := link.PacketN(2, link.BigEndianBO, link.LittleEndianBF)
	
	ms.server, err = link.Listen(cfg.TransportProtocols, cfg.Listen, protocol)
	if err != nil {
		panic(err)
	}
	log.Println("server start:", ms.server.Listener().Addr().String())
	
	ms.initChannels()

	ms.server.AcceptLoop(func(session *link.Session) {
		log.Println("client", session.Conn().RemoteAddr().String(), "in")
		
		inMsg, err := session.Read()
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println(string(inMsg.Get()))
		
		sessionStore := NewSessionStore()
		sessionStore.ClientID = string(inMsg.Get())
		sessionStore.ClientAddr = session.Conn().RemoteAddr().String()
		sessionStore.MsgServerAddr = cfg.LocalIP
		sessionStore.ID = strconv.FormatUint(session.Id(), 10)
		
		err = session.Send(link.JSON {
			sessionStore,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
	})
}
