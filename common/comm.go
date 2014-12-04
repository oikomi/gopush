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

package common

import (
	"math/rand"
	"time"
	//"github.com/oikomi/gopush/session_manager/redis_store"
)

const KeyPrefix string = "push"

var DefaultRedisConnectTimeout uint32 = 2000
var DefaultRedisReadTimeout    uint32 = 1000
var DefaultRedisWriteTimeout   uint32 = 1000

var DefaultRedisOptions RedisStoreOptions = RedisStoreOptions {
	Network        :   "tcp",
	Address        :   ":6379",
	ConnectTimeout : time.Duration(DefaultRedisConnectTimeout)*time.Millisecond,
	ReadTimeout    : time.Duration(DefaultRedisReadTimeout)*time.Millisecond,
	WriteTimeout   : time.Duration(DefaultRedisWriteTimeout)*time.Millisecond,
	Database       :  1,
	KeyPrefix      : "push",
}

func SelectServer(serverList []string, serverNum int) string {
	return serverList[rand.Intn(serverNum)]
}

func GetSessionFromCID() {
	

}