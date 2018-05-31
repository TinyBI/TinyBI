// Copyright (C)2018 by Lei Peng <pyp126@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package web

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
	"tinybi/model"
)

type Session struct {
	SessionId string
	User      model.CoreUser
	AclCodes  map[string]bool
	Expire    int64
}

type SessionManager struct {
	sessions map[string]*Session
	mutex    *sync.RWMutex
}

func NewSessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func NewSession() *Session {
	session := new(Session)
	session.SessionId = NewSessionId()
	session.AclCodes = make(map[string]bool)
	session.Expire = time.Now().Unix() + SessionTimeout
	return session
}

func NewSessionManager() *SessionManager {
	sessionManger := new(SessionManager)
	sessionManger.sessions = make(map[string]*Session)
	sessionManger.mutex = new(sync.RWMutex)
	return sessionManger
}

func (this *SessionManager) Get(sessionId string) *Session {
	if sessionId == "" {
		return nil
	}
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	session, ok := this.sessions[sessionId]
	if ok {
		return session
	}
	return nil
}

func (this *SessionManager) Set(session *Session) error {
	if session.SessionId == "" {
		return errors.New("Empty session ID")
	}
	_, ok := this.sessions[session.SessionId]
	if ok {
		errStr := fmt.Sprintf("Duplicate session ID:%s", session.SessionId)
		return errors.New(errStr)
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.sessions[session.SessionId] = session
	return nil
}

func (this *SessionManager) Delete(sessionId string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.sessions, sessionId)
}

func (this *SessionManager) Gc() {
	currentTimestamp := time.Now().Unix()
	for _, session := range this.sessions {
		if currentTimestamp >= session.Expire {
			this.Delete(session.SessionId)
		}
	}
}

func (this *SessionManager) Start() {
	go func() {
		for {
			select {
			case <-time.After(SessionGcInterval * time.Second):
				this.Gc()
				break
			}
		}
	}()
}

const SessionTimeout int64 = 900
const SessionGcInterval time.Duration = 120

var Sessions *SessionManager

func init() {
	Sessions = NewSessionManager()
	Sessions.Start()
}
