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
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"tinybi/core"
	"tinybi/model"
	"tinybi/web"
)

type ProfileApp struct {
	web.BaseWebApp
}

func NewProfileApp() *ProfileApp {
	app := new(ProfileApp)
	app.AppVersion = new(web.VersionInfo)
	app.AppVersion.Name = "Profile"
	app.AppVersion.Version = "0.0.1"
	app.AppVersion.Description = "User Profile"
	return app
}

func (this ProfileApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		switch r.URL.Query().Get("act") {
		case "user":
			this.showUserPage(w, r)
			break
		default:
			web.ErrorNotFound(w, r)
		}
	} else {
		switch r.URL.Query().Get("act") {
		case "updateProfile":
			this.updateProfile(w, r)
			break
		default:
			web.ErrorNotFound(w, r)
		}
	}
}

func (this ProfileApp) showUserPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		User model.CoreUser
		Info struct {
			Show    bool
			Type    string
			Message string
		}
	}
	lang := web.GetUILang(w, r)
	if web.AclRedirect(w, r, "INDEX", "/login.html") {
		cookie, err := r.Cookie(web.SessionCookieId)
		if err != nil {
			log.Println(err)
			web.ErrorNotFound(w, r)
			return
		}
		session := web.Sessions.Get(cookie.Value)
		if session == nil {
			web.ErrorNotFound(w, r)
			return
		}
		Html.User = session.User
		err = web.GetTemplate(lang, "core/user_profile.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	}
}

func (this ProfileApp) updateProfile(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		User model.CoreUser
		Info struct {
			Show    bool
			Type    string
			Message string
		}
	}
	lang := web.GetUILang(w, r)
	if web.AclRedirect(w, r, "INDEX", "/login.html") {
		cookie, err := r.Cookie(web.SessionCookieId)
		if err != nil {
			log.Println(err)
			web.ErrorNotFound(w, r)
			return
		}
		session := web.Sessions.Get(cookie.Value)
		if session == nil {
			web.ErrorNotFound(w, r)
			return
		}
		Html.User = session.User
		r.ParseForm()
		curPassword := r.Form.Get("password")
		//To update profile, you must enter the correct password;
		if curPassword == "" {
			Html.Info.Show = true
			Html.Info.Type = "danger"
			Html.Info.Message = "To update profile, you must enter your current password"
		} else {
			//Validate password;
			var valPass struct {
				Pass1 string
				Pass2 string
			}
			selectStr := fmt.Sprintf("md5( concat( md5('%s'), password_salt )) as pass1, `password` as pass2", curPassword)
			ok, err := core.DB.Table("core_users").Select(selectStr).Where("id=?", Html.User.Id).Get(&valPass)
			if !ok {
				log.Println("Illegal visit profile page")
				web.ErrorNotFound(w, r)
				return
			}
			if err != nil {
				log.Println(err)
				web.ErrorNotFound(w, r)
				return
			}
			if valPass.Pass1 != valPass.Pass2 {
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = "Wrong password"
			}
		}
		if !Html.Info.Show {
			newPassword := r.Form.Get("newPassword")
			if newPassword != "" {
				//Update password;
				newPassword2 := r.Form.Get("newPassword2")
				if newPassword != newPassword2 {
					Html.Info.Show = true
					Html.Info.Type = "danger"
					Html.Info.Message = "The new password does not correspond with the confirmed password"
				} else {
					pass, err := core.DB.SQL("select md5( concat( md5(?), ? )) as pass", newPassword, Html.User.PasswordSalt).QueryString()
					if err != nil {
						Html.Info.Show = true
						Html.Info.Type = "danger"
						Html.Info.Message = err.Error()
					} else {
						Html.User.LastPassword = Html.User.Password
						Html.User.LastPasswordTime = time.Now()
						Html.User.Password = pass[0]["pass"]
						Html.User.Lang = r.Form.Get("lang")
						_, err = core.DB.Table("core_users").Where("id=?", Html.User.Id).Update(&Html.User)
						if err != nil {
							Html.Info.Show = true
							Html.Info.Type = "danger"
							Html.Info.Message = err.Error()
						} else {
							//Update session info;
							session.User.LastPassword = Html.User.LastPassword
							session.User.LastPasswordTime = Html.User.LastPasswordTime
							session.User.Password = Html.User.Password
						}
					}
				}
			}
			if r.Form.Get("nickname") != Html.User.NickName {
				if Html.User.NickName == "" {
					Html.Info.Show = true
					Html.Info.Type = "danger"
					Html.Info.Message = "You must enter a valid and unique nick"
				} else {
					Html.User.NickName = r.Form.Get("nickname")
					Html.User.Lang = r.Form.Get("lang")
					_, err = core.DB.Table("core_users").Where("id=?", Html.User.Id).Update(&Html.User)
					if err != nil {
						Html.Info.Show = true
						Html.Info.Type = "danger"
						Html.Info.Message = err.Error()
					} else {
						//Update session info;
						session.User.NickName = Html.User.NickName
					}
				}
			}
		}
		if !Html.Info.Show {
			Html.Info.Show = true
			Html.Info.Type = "success"
			Html.Info.Message = "Your profile has been updated successfully"
		}
		err = web.GetTemplate(lang, "core/user_profile.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	}
}

var PModWebApp *ProfileApp

var ModWebApp ProfileApp

func init() {
	PModWebApp = NewProfileApp()
	ModWebApp = *PModWebApp
}
