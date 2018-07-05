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
package web

import (
	"encoding/json"
	"log"
	"net/http"
)

const SessionCookieId string = "ECS8086"

func AclRedirect(w http.ResponseWriter, r *http.Request, code string, url string) bool {
	cookie, err := r.Cookie(SessionCookieId)
	if err != nil {
		http.Redirect(w, r, url, http.StatusFound)
		return false
	}
	session := Sessions.Get(cookie.Value)
	if session != nil {
		//All users can visit home page;
		if code == "INDEX" {
			return true
		}
		_, ok := session.AclCodes[code]
		if ok {
			return true
		}
	}
	http.Redirect(w, r, url, http.StatusFound)
	return false
}

func GetUILang(w http.ResponseWriter, r *http.Request) string {
	//Default language;
	var lang string = "en_US"
	//Try to get language code from cookie;
	langCookie, err := r.Cookie("lang")
	if err != nil {
		//From URL parameter;
		lang = r.URL.Query().Get("lang")
		if lang == "" {
			lang = "en_US"
		}
		nLangCookie := http.Cookie{Name: "lang", Value: lang}
		http.SetCookie(w, &nLangCookie)
	} else {
		lang = langCookie.Value
	}

	return lang
}

func SetUILang(w http.ResponseWriter, r *http.Request, lang string) string {
	oldLang := GetUILang(w, r)
	nLangCookie := http.Cookie{Name: "lang", Value: lang}
	http.SetCookie(w, &nLangCookie)
	return oldLang
}

func JsonEncode(in interface{}) string {
	bstr, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	return string(bstr)
}

func ErrorNotFound(w http.ResponseWriter, r *http.Request) {
	lang := GetUILang(w, r)
	w.WriteHeader(http.StatusNotFound)
	err := GetTemplate(lang, "core/404.html").Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}
