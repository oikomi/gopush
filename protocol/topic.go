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

import (
	"github.com/funny/link"
)

type TopicMap   map[string]*Topic

type Topic struct {
	TopicName     string
	Channel       *link.Channel
	TA            *TopicAttribute
	ClientIdList  []string
}

func NewTopic(topicName string, CreaterID string, CreaterSession *link.Session) *Topic {
	return &Topic {
		TopicName    : topicName,
		Channel      : new(link.Channel),
		TA           : NewTopicAttribute(CreaterID, CreaterSession),
		ClientIdList : make([]string, 0),
	}
}

type TopicAttribute struct {
	CreaterID          string
	CreaterSession     *link.Session
}

func NewTopicAttribute(CreaterID string, CreaterSession *link.Session) *TopicAttribute {
	return &TopicAttribute {
		CreaterID      : CreaterID,
		CreaterSession : CreaterSession,
	}
}