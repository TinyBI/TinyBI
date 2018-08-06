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
	"tinybi/web"
	"tinybi/model"
	"net/http"
	"sort"
	"crypto/md5"
	"io"
	"fmt"
	"tinybi_mods/api/mailer"
)

type ApiResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ApiApp struct {
	web.BaseWebApp
}

func NewApiApp() *ApiApp {
	app := new(ApiApp)
	app.AppVersion = new(web.VersionInfo)
	app.AppVersion.Name = "API"
	app.AppVersion.Version = "0.0.1"
	app.AppVersion.Description = "WEB APIs"
	return app
}

func (this ApiApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp ApiResponse
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		resp.Status = false
		resp.Message = "Currently, the server only support GET & POST calling"
		w.Write([]byte(web.JsonEncode(resp)))
		return
	}
	appKey := r.URL.Query().Get("app_key")
	appKeyVal := model.BusinessSettings.Get("APP_KEY")
	if appKeyVal == nil || appKey == "" || appKeyVal.Value != appKey {
		resp.Status = false
		resp.Message = "Illegal call with wrong app key"
		w.Write([]byte(web.JsonEncode(resp)))
		return
	}
	appSecret := model.BusinessSettings.Get("APP_SECRET")
	if appSecret == nil || appSecret.Value == "" {
		resp.Status = false
		resp.Message = "This server disabled API calling"
		w.Write([]byte(web.JsonEncode(resp)))
		return
	}
	sign := r.URL.Query().Get("sign")
	if sign == "" {
		resp.Status = false
		resp.Message = "Illegal call with empty signature"
		w.Write([]byte(web.JsonEncode(resp)))
		return
	}
	var keys []string
	for key, _ := range r.URL.Query() {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	//Make sign;
	rawSign := appSecret.Value
	for _, key := range keys {
		if key != "sign" {
			rawSign += key
			rawSign += r.URL.Query().Get(key)
		}
	}
	rawSign += appSecret.Value
	signHash := md5.New()
	io.WriteString(signHash, rawSign)
	signVal := fmt.Sprintf("%x", signHash.Sum(nil))
	if signVal != sign {
		resp.Status = false
		resp.Message = "Illegal call with wrong signature"
		w.Write([]byte(web.JsonEncode(resp)))
		return
	}
	if r.Method == http.MethodGet {
		switch r.URL.Query().Get("service") {
		default:
			web.ErrorNotFound(w, r)
		}

	}
	if r.Method == http.MethodPost {
		switch r.URL.Query().Get("service") {
		case "mailer":
			resp.Status, resp.Message, resp.Data = mailer.Post(w, r)
			w.Write([]byte(web.JsonEncode(resp)))
			return
		default:
			web.ErrorNotFound(w, r)
		}

	}
}

var PModWebApp *ApiApp

var ModWebApp ApiApp

func init() {
	PModWebApp = NewApiApp()
	ModWebApp = *PModWebApp
}
