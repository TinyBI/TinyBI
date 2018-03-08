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
	"log"
	"net/http"
	"tinybi/core"
	"tinybi/webcore"
)

type LoginApp struct {
	webcore.BaseWebApp
}

func (this LoginApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		switch r.URL.Query().Get("act") {
		case "setLang":
			this.setUILang(w, r)
			break
		default:
			this.showForm(w, r)
		}

	} else {
		//Do login action;
		r.ParseForm()
	}
}

func (this LoginApp) setUILang(w http.ResponseWriter, r *http.Request) {
	//Default language;
	lang := "en_US"
	lang = r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en_US"
	}
	webcore.SetUILang(w, r, lang)
	http.Redirect(w, r, "/login.html", http.StatusFound)
}

func (this LoginApp) showForm(w http.ResponseWriter, r *http.Request) {
	lang := webcore.GetUILang(w, r)
	if core.Conf.Debug {
		log.Println(lang)
	}
	//Show login form;
	webcore.Render(w, lang, "login_form.html", nil)
}
