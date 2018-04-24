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
package webcore

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"tinybi/core"

	"github.com/satori/go.uuid"
)

type UserInfo struct {
	Id       int64
	UserName string
	NickName string
	RoleName string
}

type Session struct {
	SessionId string
	User      UserInfo
	AclRoles  map[string]string
	Expire    int64
}

var sessions map[string]*Session
var sessionMutex *sync.RWMutex

func init() {
	sessions = make(map[string]*Session)
	sessionMutex = new(sync.RWMutex)
}

func AclCheck(sessionId string, aclCode string) bool {
	if sessionId == "" || aclCode == "" {
		return false
	}
	session := GetSession(sessionId)
	if session == nil {
		return false
	}
	if session.Expire < time.Now().Unix() {
		if core.Conf.Debug {
			log.Printf("(%d)Session[%s] is expired at %d\n", time.Now().Unix(), sessionId, session.Expire)
		}
		RemoveSession(session)
		return false
	}
	//All active users have the right to visit index page;
	if aclCode == "INDEX" {
		return true
	}
	_, ok := session.AclRoles[aclCode]
	if !ok {
		return false
	}
	return true
}

func AclCheckRedirect(w http.ResponseWriter, r *http.Request, aclCode, rdirectUrl string) {
	sessionCookie, err := r.Cookie("session")
	sessionId := ""
	if err != nil {
		sessionId = r.URL.Query().Get("sId")
	} else {
		sessionId = sessionCookie.Value
	}
	if !AclCheck(sessionId, aclCode) {
		http.Redirect(w, r, rdirectUrl, http.StatusFound)
		return
	}
}

func AclGetSession(r *http.Request) *Session {
	sessionCookie, err := r.Cookie("session")
	sessionId := ""
	if err != nil {
		sessionId = r.URL.Query().Get("sId")
	} else {
		sessionId = sessionCookie.Value
	}
	return GetSession(sessionId)
}

func NewSession() (*Session, error) {
	session := new(Session)
	session.AclRoles = make(map[string]string)
	u2, err := uuid.NewV4()
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to generate UUID", err)
		}
		return nil, err
	}
	session.SessionId = u2.String()
	session.Expire = time.Now().Unix() + core.Conf.App.Web.SessionTimeout
	if core.Conf.Debug {
		log.Printf("Session Timeout:%d\n", core.Conf.App.Web.SessionTimeout)
		log.Printf("Session[%s] will be expired at %d\n", session.SessionId, session.Expire)
	}
	return session, nil
}

func GetSession(sessionId string) *Session {
	if sessionId == "" {
		return nil
	}
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	session, ok := sessions[sessionId]
	if ok {
		return session
	}
	return nil
}

func SetSession(session *Session) error {
	if session == nil {
		return errors.New("Null Pointer of Session")
	}
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	sessions[session.SessionId] = session
	return nil
}

func RemoveSession(session *Session) {
	if session != nil {
		sessionMutex.Lock()
		defer sessionMutex.Unlock()
		delete(sessions, session.SessionId)
	}
}

func EmailLogin(email, password string) *Session {
	if email == "" || password == "" {
		return nil
	}
	sql := "SELECT core_users.id AS user_id, core_users.user_name, "
	sql += "core_users.nick_name, core_roles.role_name, core_roles.role_codes "
	sql += "FROM core_users, core_roles WHERE core_users.role_id = core_roles.id"
	sql += " AND core_users.user_name = ? "
	sql += "AND md5( concat( md5(?), core_users.password_salt )) = core_users.`password` "
	sql += "AND core_users.`status` = 'ACTIVE'"
	row, err := core.DB.Query(sql, email, password)
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query user", err)
		}
		return nil
	}
	defer row.Close()
	session, err := NewSession()
	if err != nil {
		if core.Conf.Debug {
			log.Println("System Error", err)
		}
		return nil
	}
	for row.Next() {
		var roleCodes string
		row.Scan(&session.User.Id, &session.User.UserName, &session.User.NickName, &session.User.RoleName, &roleCodes)
		roles := strings.Split(roleCodes, ",")
		session.AclRoles = make(map[string]string)
		for _, v := range roles {
			session.AclRoles[v] = v
		}
		err = SetSession(session)
		if err != nil {
			if core.Conf.Debug {
				log.Println("System Error", err)
			}
			return nil
		}
		return session
	}
	//Normal error;
	//Wrong username or password;
	return nil
}

//Is the nickname has already be taken?
func isNicknameAvaible(userId int64, nickName string) bool {
	if userId == 0 || nickName == "" {
		return false
	}
	sql := "SELECT core_users.id AS user_id "
	sql += "FROM core_users "
	sql += "WHERE id != ? "
	sql += "AND nick_name = ?"
	row, err := core.DB.Query(sql, userId, nickName)
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query user", err)
		}
		return false
	}
	defer row.Close()
	for row.Next() {
		return false
	}
	return true
}

//true when success;
//the error message can be retrieved from the second return value;
//this method will update the password if newPassword is not empty
func UpdateMyProfile(session *Session, info UserInfo, password string, newPassword string) (bool, string) {
	if session == nil {
		return false, "Illegal call"
	}
	//Check the current password;
	sql := "SELECT core_users.id AS user_id "
	sql += "FROM core_users "
	sql += "WHERE id = ? "
	sql += "AND md5( concat( md5(?), core_users.password_salt )) = core_users.`password` "
	row, err := core.DB.Query(sql, session.User.Id, password)
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query user", err)
		}
		return false, "Wrong password"
	}
	defer row.Close()
	for row.Next() {
		var userId int64
		row.Scan(&userId)
		if userId != session.User.Id {
			return false, "You are not granted to update another one's profile"
		}
		if !isNicknameAvaible(userId, info.NickName) {
			return false, "The username has already been taken"
		}
		sql := "UPDATE core_users "
		sql += "SET nick_name = ? "
		if newPassword != "" {
			sql += ",password = md5( concat( md5(?), core_users.password_salt )) "
		}
		sql += "WHERE id = ? "
		if newPassword != "" {
			_, err = core.DB.Query(sql, info.NickName, newPassword, userId)
		} else {
			_, err = core.DB.Query(sql, info.NickName, userId)
		}
		if err != nil {
			if core.Conf.Debug {
				log.Println("Fail to update user", err)
			}
			return false, "System error, please contact administrators"
		}
		session.User.NickName = info.NickName
		err = SetSession(session)
		if err != nil {
			if core.Conf.Debug {
				log.Println("Fail to update session", err)
			}
			return false, "System error, please contact administrators"
		}
		return true, ""
	}
	return false, "Wrong password"
}

func GcSession() {
	if core.Conf.Debug {
		log.Println("Webcore::GcSession Start")
	}
	for _, session := range sessions {
		if session.Expire < time.Now().Unix() {
			if core.Conf.Debug {
				log.Println("Webcore::GcSession> Remove expired session:", session.SessionId)
			}
			RemoveSession(session)
		}
	}
	if core.Conf.Debug {
		log.Println("Webcore::GcSession End")
	}
}
