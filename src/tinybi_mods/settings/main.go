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
	"strconv"
	"tinybi/core"
	"tinybi/model"
	"tinybi/web"

	"github.com/chai2010/gettext-go/gettext"
)

type SettingsApp struct {
	web.BaseWebApp
}

func NewSettingsApp() *SettingsApp {
	app := new(SettingsApp)
	app.AppVersion = new(web.VersionInfo)
	app.AppVersion.Name = "Settings"
	app.AppVersion.Version = "0.0.1"
	app.AppVersion.Description = "Version Viewer"
	return app
}

func (this SettingsApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		switch r.URL.Query().Get("act") {
		case "list":
			this.list(w, r)
			break
		case "add":
			this.editPage("add", w, r)
			break
		case "edit":
			this.editPage("edit", w, r)
			break
		default:
			this.showPage(w, r)
		}

	} else {
		switch r.URL.Query().Get("act") {
		case "addExec":
			this.editExec("add", w, r)
			break
		case "editExec":
			this.editExec("edit", w, r)
			break
		default:
			web.ErrorNotFound(w, r)
		}
	}
}

func (this SettingsApp) showPage(w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	if web.AclRedirect(w, r, "SETTINGS_ADMIN_VIEW", "/login.html") {
		//Show Page;
		err := web.GetTemplate(lang, "settings/index.html").Execute(w, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

func (this SettingsApp) list(w http.ResponseWriter, r *http.Request) {
	if web.AclRedirect(w, r, "SETTINGS_ADMIN_VIEW", "/login.html") {
		w.Header().Set("Content-Type", "application/json")
		nullRet := `{"data":[]}`
		var fullRet struct {
			Data []struct {
				Code        string `json:"0"`
				Description string `json:"1"`
				Value       string `json:"2"`
				EditUrl     string `json:"3"`
			} `json:"data"`
		}
		rawList := model.BusinessSettings.List()
		if len(rawList) == 0 {
			w.Write([]byte(nullRet))
			return
		}
		for _, raw := range rawList {
			var fRow struct {
				Code        string `json:"0"`
				Description string `json:"1"`
				Value       string `json:"2"`
				EditUrl     string `json:"3"`
			}
			fRow.Code = raw.Code
			fRow.Description = raw.Description
			fRow.Value = raw.Value
			editI18n := gettext.Gettext("Edit")
			fRow.EditUrl = fmt.Sprintf(`<a href="/settings.html?act=edit&id=%d">%s</a>`, raw.Id, editI18n)
			fullRet.Data = append(fullRet.Data, fRow)
		}
		w.Write([]byte(web.JsonEncode(fullRet)))
	}
}

func (this SettingsApp) editPage(act string, w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	var Html struct {
		Title    string
		Act      string
		Settings model.Settings
		Info     struct {
			Show    bool
			Type    string
			Message string
		}
	}
	if web.AclRedirect(w, r, "SETTINGS_ADMIN_EDIT", "/login.html") {
		switch act {
		case "add":
			Html.Title = "New Settings"
			Html.Act = "addExec"
			break
		case "edit":
			Html.Title = "Edit Settings"
			Html.Act = "editExec"
			//Load User info;
			settingsId := r.URL.Query().Get("id")
			if settingsId == "" {
				log.Println("Visit settings edit page with null ID")
				web.ErrorNotFound(w, r)
				return
			}
			ok, err := core.DB.Table("core_settings").Where("id=?", settingsId).Get(&Html.Settings)
			if !ok {
				log.Println("Visit user settings page with illegal ID")
				if err != nil {
					log.Println(err)
				}
				web.ErrorNotFound(w, r)
				return
			}
			break
		}
		err := web.GetTemplate(lang, "settings/editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	}
}

func (this SettingsApp) editExec(act string, w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	var Html struct {
		Title    string
		Act      string
		Settings model.Settings
		Info     struct {
			Show    bool
			Type    string
			Message string
		}
	}
	if web.AclRedirect(w, r, "SETTINGS_ADMIN_EDIT", "/login.html") {
		r.ParseForm()
		Html.Settings.Code = r.Form.Get("code")
		Html.Settings.Description = r.Form.Get("description")
		Html.Settings.Value = r.Form.Get("value")
		switch act {
		case "add":
			Html.Title = "New Settings"
			Html.Act = "addExec"
			break
		case "edit":
			Html.Title = "Edit Settings"
			Html.Act = "editExec"
			settingsId := r.URL.Query().Get("id")
			if settingsId == "" {
				log.Println("Visit settings edit execution page with null ID")
				web.ErrorNotFound(w, r)
				return
			}
			sId, err := strconv.Atoi(settingsId)
			if err != nil {
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = err.Error()
			}
			Html.Settings.Id = int64(sId)
			break
		}
		if !Html.Info.Show {
			err := model.BusinessSettings.Set(&Html.Settings)
			if err != nil {
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = err.Error()
			}
		}
		if Html.Info.Show {
			err := web.GetTemplate(lang, "settings/editor.html").Execute(w, Html)
			if err != nil {
				log.Println(err)
			}
		} else {
			http.Redirect(w, r, "/settings.html", http.StatusFound)
		}
	}
}

var PModWebApp *SettingsApp

var ModWebApp SettingsApp

func init() {
	PModWebApp = NewSettingsApp()
	ModWebApp = *PModWebApp
}
