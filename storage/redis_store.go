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

package storage

import (
	"time"
	"sync"
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
)

var (
	ErrNoKeyPrefix = errors.New("cannot get session keys without a key prefix")
)

type StoreSession struct {
	ClientID        string
	ClientAddr      string
	MsgServerAddr   string
	ID              string
	MaxAge          time.Duration
}

type RedisStoreOptions struct {
	Network              string
	Address              string
	ConnectTimeout       time.Duration
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	Database             int           // Redis database to use for session keys
	KeyPrefix            string        // If set, keys will be KeyPrefix:SessionID (semicolon added)
	BrowserSessServerTTL time.Duration // Defaults to 2 days
}

type RedisStore struct {
	opts        *RedisStoreOptions
	conn        redis.Conn
	rwMutex     sync.Mutex
}

// Create a redis session store with the specified options.
func NewRedisStore(opts *RedisStoreOptions) *RedisStore {
	var err error
	rs := &RedisStore{
		opts : opts, 
		conn : nil,
		}
	rs.conn, err = redis.DialTimeout(opts.Network, opts.Address, opts.ConnectTimeout,
		opts.ReadTimeout, opts.WriteTimeout)
	if err != nil {
		panic(err)
	}
	return rs
}

// Get the session from the store.
func (self *RedisStore) Get(id string) (*StoreSession, error) {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	key := id
	if self.opts.KeyPrefix != "" {
		key = self.opts.KeyPrefix + ":" + id
	}
	b, err := redis.Bytes(self.conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	var sess StoreSession
	err = json.Unmarshal(b, &sess)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// Save the session into the store.
func (self *RedisStore) Set(sess *StoreSession) error {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	b, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	key := sess.ClientID
	if self.opts.KeyPrefix != "" {
		key = self.opts.KeyPrefix + ":" + sess.ClientID
	}
	ttl := sess.MaxAge
	if ttl == 0 {
		// Browser session, set to specified TTL
		ttl = self.opts.BrowserSessServerTTL
		if ttl == 0 {
			ttl = 2 * 24 * time.Hour // Default to 2 days
		}
	}
	_, err = self.conn.Do("SETEX", key, int(ttl.Seconds()), b)
	if err != nil {
		return err
	}
	return nil
}
// Delete the session from the store.
func (self *RedisStore) Delete(id string) error {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	key := id
	if self.opts.KeyPrefix != "" {
		key = self.opts.KeyPrefix + ":" + id
	}
	_, err := self.conn.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}
// Clear all sessions from the store. Requires the use of a key
// prefix in the store options, otherwise the method refuses to delete all keys.
func (self *RedisStore) Clear() error {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	vals, err := self.getSessionKeys()
	if err != nil {
		return err
	}
	if len(vals) > 0 {
		self.conn.Send("MULTI")
		for _, v := range vals {
			self.conn.Send("DEL", v)
		}
		_, err = self.conn.Do("EXEC")
		if err != nil {
			return err
		}
	}
	return nil
}
// Get the number of session keys in the store. Requires the use of a
// key prefix in the store options, otherwise returns -1 (cannot tell
// session keys from other keys).
func (self *RedisStore) Len() int {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	vals, err := self.getSessionKeys()
	if err != nil {
		return -1
	}
	return len(vals)
}
func (self *RedisStore) getSessionKeys() ([]interface{}, error) {
	self.rwMutex.Lock()
	defer self.rwMutex.Unlock()
	if self.opts.KeyPrefix != "" {
		return redis.Values(self.conn.Do("KEYS", self.opts.KeyPrefix+":*"))
	}
	return nil, ErrNoKeyPrefix
}

