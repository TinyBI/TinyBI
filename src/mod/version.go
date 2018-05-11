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
	"tinybi/webcore"
)

// This is the demo module for TinyBI;
type modVersion struct {
	webcore.BaseWebApp
}

func (this modVersion) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		webcore.AclCheckRedirect(w, r, "INDEX", "/login.html")
		switch r.URL.Query().Get("act") {
		default:
			//Show WEB page;
			this.showPage(w, r)
			break
		}
	} else {
		//This app does not accept POST action;
	}
}

func (this modVersion) showPage(w http.ResponseWriter, r *http.Request) {
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "version/index.html").Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

var ModVersion modVersion
