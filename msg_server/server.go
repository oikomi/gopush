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
	//"fmt"
	"time"
	"flag"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/funny/link"
	"github.com/oikomi/gopush/protocol"
	"github.com/oikomi/gopush/storage"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

type MsgServer struct {
	cfg         *MsgServerConfig
	sessions    SessionMap
	channels    ChannelMap
	server      *link.Server
	redisStore  *storage.RedisStore
}

func NewMsgServer(cfg *MsgServerConfig) *MsgServer {
	return &MsgServer {
		cfg        : cfg,
		sessions   : make(SessionMap),
		channels   : make(ChannelMap),
		server     : new(link.Server),
		redisStore : storage.NewRedisStore(&storage.RedisStoreOptions {
			Network :   "tcp",
			Address :   cfg.Redis.Port,
			ConnectTimeout : time.Duration(cfg.Redis.ConnectTimeout)*time.Millisecond,
			ReadTimeout : time.Duration(cfg.Redis.ReadTimeout)*time.Millisecond,
			WriteTimeout : time.Duration(cfg.Redis.WriteTimeout)*time.Millisecond,
			Database :  1,
			KeyPrefix : "push",
		}),
	}
}

func (self *MsgServer)createChannels() {
	for _, c := range ChannleList {
		channel := link.NewChannel(self.server.Protocol())
		self.channels[c] = channel
	}
}

func (self *MsgServer)parseProtocol(cmd []byte, session *link.Session) error {
	var c protocol.Cmd
	
	err := json.Unmarshal(cmd, &c)
	if err != nil {
		glog.Error("error:", err)
		return err
	}
	
	pp := NewProtoProc(self)
	
	glog.Info(c.CmdName)

	switch c.CmdName {
		case protocol.SUBSCRIBE_CHANNEL_CMD:
			pp.procSubscribeChannel(c, session)
		case protocol.SEND_CLIENT_ID_CMD:
			err = pp.procClientID(c, session)
			if err != nil {
				glog.Error("error:", err)
				return err
			}
		case protocol.SEND_MESSAGE_P2P_CMD:
			pp.procSendMessageP2P(c, session)
		}

	return err
}
