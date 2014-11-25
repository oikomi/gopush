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
	"net"
	"fmt"
	"log"
	"math/rand"
	"encoding/json"
	"github.com/oikomi/gopush/protocol"
)

func selectServer(serverList []string, serverNum int) string {
	return serverList[rand.Intn(serverNum)]
}

func getHostIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		fmt.Println(addr.String())
	}
}

func parseCmd(cmd []byte) {
	var c protocol.Cmd
	
	err := json.Unmarshal(cmd, &c)
	if err != nil {
		log.Fatalln("error:", err)
	}
	
	switch c.CmdName {
		case protocol.SUBSCRIBE_CHANNEL_CMD:
			fmt.Println("one")
		case protocol.SEND_CLIENT_ID_CMD:
			fmt.Println("two")
		}

}