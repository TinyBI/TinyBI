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
package apps

import (
	"encoding/base64"
	"log"
	"net/http"
	"tinybi/webcore"
)

type UserProfileApp struct {
	webcore.BaseWebApp
}

func (this UserProfileApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		webcore.AclCheckRedirect(w, r, "INDEX", "/login.html")
		//Show WEB page;
		this.showPage(w, r)
	} else {
		//Update Profile / Password;
		webcore.AclCheckRedirect(w, r, "INDEX", "/login.html")
		this.updateProfile(w, r)
	}
}

func (this UserProfileApp) showPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Session *webcore.Session
		Info    struct {
			Show    bool
			Type    string
			Message string
		}
	}
	switch r.URL.Query().Get("eNo") {
	case "1":
		Html.Info.Show = true
		Html.Info.Type = "danger"
		Html.Info.Message = "You must enter the current password to update your profile"
		break
	case "2":
		Html.Info.Show = true
		Html.Info.Type = "danger"
		Html.Info.Message = "The new password does not correspond with the confirmed password"
		break
	case "99":
		Html.Info.Show = true
		Html.Info.Type = "danger"
		eMsg, err := base64.StdEncoding.DecodeString(r.URL.Query().Get("eMsg"))
		if err != nil {
			eMsg = []byte("Unknow Error")
		}
		Html.Info.Message = string(eMsg)
		break
	case "200":
		Html.Info.Show = true
		Html.Info.Type = "success"
		Html.Info.Message = "Your profile has been updated successfully"
		break
	default:
		Html.Info.Show = false
		Html.Info.Message = "No Error"
		break
	}
	Html.Session = webcore.AclGetSession(r)
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "user_profile.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this UserProfileApp) updateProfile(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("sId")
	session := webcore.GetSession(sessionId)
	rUrl := "/userProfile.html"
	if session != nil {
		user := webcore.UserInfo{Id: session.User.Id,
			UserName: session.User.UserName}
		r.ParseForm()
		user.NickName = r.Form.Get("nickname")
		password := r.Form.Get("password")
		if password == "" {
			rUrl = "/userProfile.html?eNo=1"
			http.Redirect(w, r, rUrl, http.StatusFound)
			return
		}
		newPassword := r.Form.Get("newPassword")
		newPassword2 := r.Form.Get("newPassword2")
		if newPassword != newPassword2 {
			rUrl = "/userProfile.html?eNo=2"
			http.Redirect(w, r, rUrl, http.StatusFound)
			return
		}
		updated, eStr := webcore.UpdateMyProfile(session, user, password, newPassword)
		if !updated {
			eMsg := base64.StdEncoding.EncodeToString([]byte(eStr))
			rUrl = "/userProfile.html?eNo=99&eMsg="
			rUrl += eMsg
		} else {
			rUrl = "/userProfile.html?eNo=200"
		}
	}
	http.Redirect(w, r, rUrl, http.StatusFound)
}
