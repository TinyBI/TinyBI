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
	"log"
	"net/http"
	"strings"
	"tinybi/core"
	"tinybi/model"
	"tinybi/web"
)

type LoginApp struct {
	web.BaseWebApp
}

func NewLoginApp() *LoginApp {
	app := new(LoginApp)
	app.AppVersion = new(web.VersionInfo)
	app.AppVersion.Name = "Login"
	app.AppVersion.Version = "0.0.1"
	app.AppVersion.Description = "User login/logout module"
	return app
}

func (this LoginApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		switch r.URL.Query().Get("act") {
		case "setLang":
			this.setUILang(w, r)
			break
		case "exit":
			this.exit(w, r)
			break
		default:
			this.showForm(w, r)
		}

	} else {
		//Do login action;
		this.loginAction(w, r)
	}
}

func (this LoginApp) setUILang(w http.ResponseWriter, r *http.Request) {
	//Default language;
	lang := "en_US"
	lang = r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en_US"
	}
	web.SetUILang(w, r, lang)
	http.Redirect(w, r, "/login.html", http.StatusFound)
}

func (this LoginApp) exit(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session")
	if err == nil {
		session := web.Sessions.Get(sessionCookie.Value)
		if session != nil {
			web.Sessions.Delete(session.SessionId)
		}
		sessionCookie.Value = ""
		sessionCookie.MaxAge = 0
		http.SetCookie(w, sessionCookie)
	}
	http.Redirect(w, r, "/login.html", http.StatusFound)
}

func (this LoginApp) showForm(w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	if core.Conf.Debug {
		log.Println(lang)
	}
	var Info struct {
		Show    bool
		Message string
	}

	switch r.URL.Query().Get("eNo") {
	case "1":
		Info.Show = true
		Info.Message = "Wrong Email address or password"
		break
	default:
		Info.Show = false
		Info.Message = "No Error"
		break
	}
	//Show login form;
	err := web.GetTemplate(lang, "core/login_form.html").Execute(w, Info)
	if err != nil {
		log.Println(err)
	}
}

func (this LoginApp) loginAction(w http.ResponseWriter, r *http.Request) {
	redirectUrl := "/index.html"
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	session := web.NewSession()
	ok, err := core.DB.Table("core_users").Where("email=?", email).And(
		"md5( concat( md5(?), password_salt )) = `password`", password).And(
		"status='ACTIVE'").Get(&session.User)
	if ok {
		var role model.CoreRole
		ok, err = core.DB.Table("core_roles").Where("id=?", session.User.RoleId).Get(&role)
		if ok {
			roleArray := strings.Split(role.RoleCodes, ",")
			for _, roleCode := range roleArray {
				session.AclCodes[roleCode] = true
			}
			web.Sessions.Set(session)
			sessionCookie := http.Cookie{Name: web.SessionCookieId, Value: session.SessionId}
			http.SetCookie(w, &sessionCookie)
		} else {
			log.Println(err)
			redirectUrl = "/login.html?eNo=1"
		}
	} else {
		log.Println(err)
		redirectUrl = "/login.html?eNo=1"
	}
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}

var PModWebApp *LoginApp

var ModWebApp LoginApp

func init() {
	PModWebApp = NewLoginApp()
	ModWebApp = *PModWebApp
}
