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
	"encoding/json"
	"time"
	"github.com/funny/link"
	"github.com/oikomi/gopush/protocol"
	"github.com/oikomi/gopush/storage"
)

type Router struct {
	cfg         *RouterConfig
	redisStore  *storage.RedisStore
}   

func NewRouter(cfg *RouterConfig) *Router {
	return &Router {
		cfg : cfg,
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

func (self *Router)connectMsgServer(ms string) (*link.Session, error) {
	p := link.PacketN(2, link.BigEndianBO, link.LittleEndianBF)
	client, err := link.Dial("tcp", ms, p)
	if err != nil {
		glog.Error(err.Error())
		panic(err)
	}

	return client, err
}

func (self *Router)handleMsgServerClient(msc *link.Session) {
	msc.ReadLoop(func(msg link.InBuffer) {
		glog.Info("msg_server", msc.Conn().RemoteAddr().String()," say: ", string(msg.Get()))
		var c protocol.Cmd
		pp := NewProtoProc(self)
		err := json.Unmarshal(msg.Get(), &c)
		if err != nil {
			glog.Error("error:", err)
		}
		switch c.CmdName {
			case protocol.SEND_MESSAGE_P2P_CMD:
				pp.procSendMsgP2P(c, msc)
				
			}
	})
}

func (self *Router)subscribeChannels() error {
	glog.Info("subscribeChannels")
	var msgServerClientList []*link.Session
	for _, ms := range self.cfg.MsgServerList {
		msgServerClient, err := self.connectMsgServer(ms)
		if err != nil {
			glog.Error(err.Error())
			return err
		}
		cmd := protocol.NewCmd()
		
		cmd.CmdName = protocol.SUBSCRIBE_CHANNEL_CMD
		cmd.Args = append(cmd.Args, protocol.SYSCTRL_SEND)
		
		err = msgServerClient.Send(link.JSON {
			cmd,
		})
		if err != nil {
			glog.Error(err.Error())
			return err
		}
		
		msgServerClientList = append(msgServerClientList, msgServerClient)
	}

	for _, msc := range msgServerClientList {
		go self.handleMsgServerClient(msc)
	}
	return nil
}
