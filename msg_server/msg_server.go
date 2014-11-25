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
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/funny/link"
	"github.com/oikomi/gopush/protocol"
)

var InputConfFile = flag.String("conf_file", "msg_server.json", "input conf file name")   

type MsgServer struct {
	cfg Config
	sessions SessionMap
	channels ChannelMap
	server *link.Server
}

func NewMsgServer() *MsgServer {
	ms := &MsgServer {
		sessions : make(SessionMap),
		channels : make(ChannelMap),
		server : new(link.Server),
	}
	
	return ms
}

func (self *MsgServer)initChannels() {
	channel := link.NewChannel(self.server.Protocol())
	self.channels[SYSCTRL_CLIENT_STATUS] = channel
}

func (self *MsgServer)procClientID(cmd protocol.Cmd, session *link.Session) {
	
	sessionStore := NewSessionStore()
	sessionStore.ClientID = string(cmd.Args[0])
	sessionStore.ClientAddr = session.Conn().RemoteAddr().String()
	sessionStore.MsgServerAddr = self.cfg.LocalIP
	sessionStore.ID = strconv.FormatUint(session.Id(), 10)
	
	err := session.Send(link.JSON {
		sessionStore,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	self.sessions[string(cmd.Args[0])] = session
}

func (self *MsgServer)parseProtocol(cmd []byte, session *link.Session) {
	var c protocol.Cmd
	
	err := json.Unmarshal(cmd, &c)
	if err != nil {
		log.Fatalln("error:", err)
	}
	
	switch c.CmdName {
		case protocol.SUBSCRIBE_CHANNEL_CMD:
			fmt.Println("one")
		case protocol.SEND_CLIENT_ID_CMD:
			self.procClientID(c, session)
		}
}

func main() {
	flag.Parse()
	cfg, err := LoadConfig(*InputConfFile)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	
	ms := NewMsgServer()
	ms.cfg = cfg
	
	p := link.PacketN(2, link.BigEndianBO, link.LittleEndianBF)
	
	ms.server, err = link.Listen(cfg.TransportProtocols, cfg.Listen, p)
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
		
		ms.parseProtocol(inMsg.Get(), session)
	})
}
