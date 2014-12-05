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
	"github.com/oikomi/gopush/common"
)

func init() {
	flag.Set("alsologtostderr", "true")
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

func (self *ProtoProc)procSendMessageP2P(cmd protocol.Cmd, session *link.Session) error {
	glog.Info("procSendMessageP2P")
	send2ID := string(cmd.Args[0])
	send2Msg := string(cmd.Args[1])
	store_session, err := common.GetSessionFromCID(self.msgServer.redisStore, send2ID)
	if err != nil {
		glog.Warningf("no ID : %s", send2ID)
		
		return err
	}
	
	if store_session.MsgServerAddr == self.msgServer.cfg.LocalIP {
		glog.Info("in the same server")
		resp := protocol.NewCmd()
		resp.CmdName = protocol.RESP_MESSAGE_P2P_CMD
		resp.Args = append(resp.Args, send2Msg)
		
		self.msgServer.sessions[send2ID].Send(link.JSON {
			resp,
		})
		if err != nil {
			glog.Fatalln(err.Error())
		}
	}
	
	return nil
}

func (self *ProtoProc)procSubscribeChannel(cmd protocol.Cmd, session *link.Session) {
	glog.Info("procSubscribeChannel")
	channelName := string(cmd.Args[0])
	self.msgServer.channels[channelName].Join(session, nil)
}
