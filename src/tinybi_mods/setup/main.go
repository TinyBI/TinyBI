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
	"tinybi/core"
	"tinybi/model"
	"tinybi/web"
)

type SetupApp struct {
	web.BaseWebApp
}

func NewSetupApp() *SetupApp {
	app := new(SetupApp)
	app.AppVersion = new(web.VersionInfo)
	app.AppVersion.Name = "Setup"
	app.AppVersion.Version = "0.0.1"
	app.AppVersion.Description = "System Setup"
	return app
}

func (this SetupApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	if !this.isFreshSystem() {
		web.ErrorNotFound(w, r)
	}
	if r.Method == http.MethodGet {
		switch r.URL.Query().Get("act") {
		case "step1":
			this.showStep1(w, r)
			break
		case "done":
			this.showDone(w, r)
			break
		default:
			this.showLicense(w, r)
		}
	} else {
		switch r.URL.Query().Get("act") {
		case "step1":
			this.doStep1(w, r)
			break
		default:

		}
	}

}

func (this SetupApp) showLicense(w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	err := web.GetTemplate(lang, "setup/license.html").Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func (this SetupApp) showStep1(w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	err := web.GetTemplate(lang, "setup/step1.html").Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func (this SetupApp) doStep1(w http.ResponseWriter, r *http.Request) {
	role := new(model.CoreRole)
	user := new(model.CoreUser)
	role.RoleCodes = ""
	_, err := core.DB.Table("core_acl_codes").Select("group_concat(code) AS codes").Where("1=1").Get(&role.RoleCodes)
	if err != nil {
		log.Println(err)
		lang := web.GetUILang(w, r)
		err = web.GetTemplate(lang, "setup/error.html").Execute(w, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	role.RoleName = "Administrators"
	user.RoleId, err = core.DB.Table("core_roles").Insert(role)
	if err != nil {
		log.Println(err)
		lang := web.GetUILang(w, r)
		err = web.GetTemplate(lang, "setup/error.html").Execute(w, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	user.RoleId = role.Id
	r.ParseForm()
	user.EMail = r.Form.Get("email")
	user.NickName = r.Form.Get("nickname")
	user.Password = r.Form.Get("password")
	user.PasswordSalt = core.RandomString(64)
	user.Password = model.CoreUserPassword(user.Password, user.PasswordSalt)
	user.Source = "A"
	user.Status = "ACTIVE"
	_, err = core.DB.Table("core_users").Insert(user)
	if err != nil {
		log.Println(err)
		lang := web.GetUILang(w, r)
		err = web.GetTemplate(lang, "setup/error.html").Execute(w, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	lang := web.GetUILang(w, r)
	err = web.GetTemplate(lang, "setup/done.html").Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func (this SetupApp) showDone(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Versions []web.VersionInfo
	}
	lang := web.GetUILang(w, r)
	err := web.GetTemplate(lang, "setup/done.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this SetupApp) isFreshSystem() bool {
	user := new(model.CoreUser)
	total, _ := core.DB.Table("core_users").Where("1=1").Count(user)
	if total > 0 {
		return false
	}
	total, _ = core.DB.Table("core_roles").Where("1=1").Count(user)
	if total > 0 {
		return false
	}
	return true
}

var PModWebApp *SetupApp

var ModWebApp SetupApp

func init() {
	PModWebApp = NewSetupApp()
	ModWebApp = *PModWebApp
}
