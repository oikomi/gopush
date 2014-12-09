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

package protocol

const (
	SEND_CLIENT_ID_CMD      = "SEND_CLIENT_ID"
	SUBSCRIBE_CHANNEL_CMD   = "SUBSCRIBE_CHANNEL"
	SEND_MESSAGE_P2P_CMD    = "SEND_MESSAGE_P2P"
	RESP_MESSAGE_P2P_CMD    = "RESP_MESSAGE_P2P"
	ROUTE_MESSAGE_P2P_CMD   = "ROUTE_MESSAGE_P2P"
	
	CREATE_TOPIC_CMD       = "CREATE_TOPIC"
)

type Cmd struct {
	CmdName string
	Args []string
}

func NewCmd() *Cmd {
	return &Cmd {
		CmdName : "",
		Args : make([]string, 0),
	}
}

func (self *Cmd)ParseCmd(msglist []string) {
	self.CmdName = msglist[1]
	self.Args = msglist[2:]
}

type ClientIDCmd struct {
	CmdName string
	ClientID string
}

type SendMessageP2PCmd struct {
	CmdName string
	ID string
	Msg string
}
