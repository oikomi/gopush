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
	"github.com/golang/glog"
	"time"
	"fmt"
	"github.com/funny/link"
	"github.com/oikomi/gopush/storage"
)

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
const char* build_time(void) {
	static const char* psz_build_time = "["__DATE__ " " __TIME__ "]";
	return psz_build_time;
}
*/
import "C"

var (
	buildTime = C.GoString(C.build_time())
)

func BuildTime() string {
	return buildTime
}

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

const VERSION string = "0.10"

func version() {
	fmt.Printf("session_manager version %s Copyright (c) 2014 Harold Miao (miaohonghit@gmail.com)  \n", VERSION)
}

var InputConfFile = flag.String("conf_file", "session_manager.json", "input conf file name")

func main() {
	version()
	fmt.Printf("built on %s\n", BuildTime())
	flag.Parse()
	cfg := NewSessionManagerConfig(*InputConfFile)
	err := cfg.LoadConfig()
	if err != nil {
		glog.Error(err.Error())
		return
	}
	
	p := link.PacketN(2, link.BigEndianBO, link.LittleEndianBF)
	
	server, err := link.Listen(cfg.TransportProtocols, cfg.Listen, p)
	if err != nil {
		glog.Error(err.Error())
	}
	glog.Info("server start:", server.Listener().Addr().String())
	
	redisOptions := storage.RedisStoreOptions {
			Network :   "tcp",
			Address :   cfg.Redis.Port,
			ConnectTimeout : time.Duration(cfg.Redis.ConnectTimeout)*time.Millisecond,
			ReadTimeout : time.Duration(cfg.Redis.ReadTimeout)*time.Millisecond,
			WriteTimeout : time.Duration(cfg.Redis.WriteTimeout)*time.Millisecond,
			Database :  1,
			KeyPrefix : "push",
	}
	
	sm := NewSessionManager(cfg)

	redisStore := storage.NewRedisStore(&redisOptions)

	go sm.subscribeChannels(redisStore)
	
	server.AcceptLoop(func(session *link.Session) {
	
	})
}
