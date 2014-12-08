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
	Router   *Router
}

func NewProtoProc(r *Router) *ProtoProc {
	return &ProtoProc {
		Router : r,
	}
}

func (self *ProtoProc)procSendMsgP2P(cmd protocol.Cmd, session *link.Session) error {
	glog.Info("procSendMsgP2P")
	var err error
	send2ID := string(cmd.Args[0])
	send2Msg := string(cmd.Args[1])
	glog.Info(send2Msg)
	store_session, err := common.GetSessionFromCID(self.Router.redisStore, send2ID)
	if err != nil {
		glog.Warningf("no ID : %s", send2ID)
		
		return err
	}
	glog.Info(store_session)
	return nil
}

