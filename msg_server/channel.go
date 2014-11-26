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
	"github.com/funny/link"
)

type ChannelMap map[string]*link.Channel
type SessionMap map[string]*link.Session

const (
	SYSCTRL_CLIENT_STATUS = "/sysctrl/client-status"
	SYSCTRL_SEND = "/sysctrl/send"
	/*
	/sysctrl/publish
	/sysctrl/send
	/sysctrl/batch-send
	/sysctrl/topic-control
	/sysctrl/topic-status 
	*/
)

var ChannleList []string

func init() {
	ChannleList = []string{SYSCTRL_CLIENT_STATUS, SYSCTRL_SEND}
}


