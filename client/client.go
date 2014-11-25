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
	"fmt"
	"flag"
	"log"
	"github.com/funny/link"
	"github.com/oikomi/gopush/protocol"
)

var InputConfFile = flag.String("conf_file", "client.json", "input conf file name")   


func main() {
	flag.Parse()
	cfg, err := LoadConfig(*InputConfFile)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	
	p := link.PacketN(2, link.BigEndianBO, link.LittleEndianBF)

	gatewayClient, err := link.Dial("tcp", cfg.GatewayServer, p)
	if err != nil {
		panic(err)
	}
	
	var input string
	if _, err := fmt.Scanf("%s\n", &input); err != nil {
		log.Fatal(err.Error())
	}
	
	inMsg, err := gatewayClient.Read()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(string(inMsg.Get()))

	gatewayClient.Close(nil)

	msgServerClient, err := link.Dial("tcp", string(inMsg.Get()), p)
	if err != nil {
		panic(err)
	}
	
	cmd := protocol.NewCmd()
	
	cmd.CmdName = protocol.SEND_CLIENT_ID_CMD
	cmd.Args = append(cmd.Args, input)
	
	err = msgServerClient.Send(link.JSON {
		cmd,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	
	defer msgServerClient.Close(nil)
	
	msgServerClient.ReadLoop(func(msg link.InBuffer) {
		log.Println("client", msgServerClient.Conn().RemoteAddr().String(),"say:", string(msg.Get()))
		
	})
}