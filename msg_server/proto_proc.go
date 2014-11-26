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
	"github.com/golang/glog"
	"strconv"
	"flag"
	"github.com/funny/link"
	"github.com/oikomi/gopush/protocol"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "3")
	flag.Set("log_dir", "false")
}

type ProtoProc struct {
	msgServer    *MsgServer
}

func NewProtoProc(msgServer *MsgServer) *ProtoProc {
	pp := &ProtoProc {
		msgServer : msgServer,
	}
	
	return pp
}

func (self *ProtoProc)procClientID(cmd protocol.Cmd, session *link.Session) error {
	glog.Info("procClientID")
	var err error
	sessionStore := NewSessionStore()
	sessionStore.ClientID = string(cmd.Args[0])
	sessionStore.ClientAddr = session.Conn().RemoteAddr().String()
	sessionStore.MsgServerAddr = self.msgServer.cfg.LocalIP
	sessionStore.ID = strconv.FormatUint(session.Id(), 10)
	
	if self.msgServer.channels[SYSCTRL_CLIENT_STATUS] != nil {
		err = self.msgServer.channels[SYSCTRL_CLIENT_STATUS].Broadcast(link.JSON {
			sessionStore,
		})
	}

	if err != nil {
		glog.Error(err.Error())
		return err
	}
	self.msgServer.sessions[string(cmd.Args[0])] = session
	
	return err
}

func (self *ProtoProc)procSubscribeChannel(cmd protocol.Cmd, session *link.Session) {
	glog.Info("procSubscribeChannel")
	channelName := string(cmd.Args[0])
	self.msgServer.channels[channelName].Join(session, nil)
}
