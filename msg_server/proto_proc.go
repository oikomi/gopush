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
	"github.com/oikomi/gopush/base"
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
	return &ProtoProc {
		msgServer : msgServer,
	}
}

func (self *ProtoProc)procPing(cmd protocol.Cmd, session *link.Session) error {
	glog.Info("procPing")
	cid := session.State.(*base.SessionState).ClientID
	self.msgServer.scanSessionMutex.Lock()
	defer self.msgServer.scanSessionMutex.Unlock()
	self.msgServer.sessions[cid].State.(*base.SessionState).Alive = true
	
	return nil
}

func (self *ProtoProc)procClientID(cmd protocol.Cmd, session *link.Session) error {
	glog.Info("procClientID")
	var err error
	sessionStore := NewSessionStore()
	sessionStore.ClientID = string(cmd.Args[0])
	sessionStore.ClientAddr = session.Conn().RemoteAddr().String()
	sessionStore.MsgServerAddr = self.msgServer.cfg.LocalIP
	sessionStore.ID = strconv.FormatUint(session.Id(), 10)
	
	if self.msgServer.channels[protocol.SYSCTRL_CLIENT_STATUS] != nil {
		err = self.msgServer.channels[protocol.SYSCTRL_CLIENT_STATUS].Broadcast(link.JSON {
			sessionStore,
		})
		if err != nil {
			glog.Error(err.Error())
			return err
		}
	}

	self.msgServer.sessions[string(cmd.Args[0])] = session
	self.msgServer.sessions[string(cmd.Args[0])].State = base.NewSessionState(true, string(cmd.Args[0]))
	
	return nil
}

func (self *ProtoProc)procSendMessageP2P(cmd protocol.Cmd, session *link.Session) error {
	glog.Info("procSendMessageP2P")
	var err error
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
		
		if self.msgServer.sessions[send2ID] != nil {
			self.msgServer.sessions[send2ID].Send(link.JSON {
				resp,
			})
			if err != nil {
				glog.Fatalln(err.Error())
			}
		}
	} else {
		if self.msgServer.channels[protocol.SYSCTRL_SEND] != nil {
			err = self.msgServer.channels[protocol.SYSCTRL_SEND].Broadcast(link.JSON {
				cmd,
			})
		}

		if err != nil {
			glog.Error(err.Error())
			return err
		}
	}
	
	return nil
}

func (self *ProtoProc)procRouteMessageP2P(cmd protocol.Cmd, session *link.Session) error {
	glog.Info("procRouteMessageP2P")
	var err error
	send2ID := string(cmd.Args[0])
	send2Msg := string(cmd.Args[1])
	_, err = common.GetSessionFromCID(self.msgServer.redisStore, send2ID)
	if err != nil {
		glog.Warningf("no ID : %s", send2ID)
		
		return err
	}

	resp := protocol.NewCmd()
	resp.CmdName = protocol.RESP_MESSAGE_P2P_CMD
	resp.Args = append(resp.Args, send2Msg)
	
	if self.msgServer.sessions[send2ID] != nil {
		self.msgServer.sessions[send2ID].Send(link.JSON {
			resp,
		})
		if err != nil {
			glog.Fatalln(err.Error())
		}
	}

	return nil
}


func (self *ProtoProc)procSendMessageTopic(cmd protocol.Cmd, session *link.Session) error {
	glog.Info("procSendMessageTopic")
	topicName := string(cmd.Args[0])
	send2Msg := string(cmd.Args[1])
	glog.Info(send2Msg)
	
	if self.msgServer.topics[topicName] != nil {
		glog.Info("topic in local server")
	} else {
		
	
	}
	
	return nil
}

func (self *ProtoProc)procSubscribeChannel(cmd protocol.Cmd, session *link.Session) {
	glog.Info("procSubscribeChannel")
	channelName := string(cmd.Args[0])
	glog.Info(channelName)
	if self.msgServer.channels[channelName] != nil {
		self.msgServer.channels[channelName].Join(session, nil)
	} else {
		glog.Warning(channelName + " is not exist")
	}
}

func (self *ProtoProc)procCreateTopic(cmd protocol.Cmd, session *link.Session) {
	glog.Info("procCreateTopic")
	topicName := string(cmd.Args[0])
	topic := protocol.NewTopic(topicName, (session.State).(*base.SessionState).ClientID, session)
	glog.Info(topic)
	topic.Channel = link.NewChannel(self.msgServer.server.Protocol())
	self.msgServer.topics[topicName] = topic
}

func (self *ProtoProc)procJoinTopic(cmd protocol.Cmd, session *link.Session) {
	glog.Info("procJoinTopic")
	topicName := string(cmd.Args[0])
	self.msgServer.topics[topicName].Channel.Join(session, nil)
}
