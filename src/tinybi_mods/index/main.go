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
	coreTask "tinybi/task"
	"tinybi/web"
	"tinybi_mods/index/task"
)

type IndexApp struct {
	web.BaseWebApp
}

func NewIndexApp() *IndexApp {
	app := new(IndexApp)
	app.AppVersion = new(web.VersionInfo)
	app.AppVersion.Name = "Index"
	app.AppVersion.Version = "0.0.1"
	app.AppVersion.Description = "Welcome page"
	return app
}

func (this IndexApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	this.showPage(w, r)
}

func (this IndexApp) showPage(w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	if web.AclRedirect(w, r, "INDEX", "/login.html") {
		//Show Page;
		err := web.GetTemplate(lang, "core/index.html").Execute(w, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

var PModWebApp *IndexApp

var ModWebApp IndexApp

func init() {
	PModWebApp = NewIndexApp()
	ModWebApp = *PModWebApp
	coreTask.RegTasks["EXAMPLE"] = task.NewExampleHandler()
}
