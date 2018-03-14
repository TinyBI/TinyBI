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
package webcore

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"tinybi/core"

	"github.com/chai2010/gettext-go/gettext"
)

type templateCache struct {
	Template *template.Template
	Expire   int64
}

var cache map[string]templateCache
var layoutsList []string

func getText(input string) string {
	return gettext.PGettext("", input)
}

func GetTemplate(w http.ResponseWriter, locale string, name string) *template.Template {
	tmpName := filepath.Base(name)
	temp, ok := cache[tmpName]
	if !ok || time.Now().Unix() > temp.Expire {
		funcMap := template.FuncMap{
			"gettext": getText,
		}
		//Load template with layouts from disk;
		pName := make([]string, 0)
		pName2 := strings.Split(name, "/")
		pName = append(pName, core.Conf.App.Web.TemplatesPath)
		pName = append(pName, pName2...)
		tFPath := filepath.Join(pName...)
		filePaths := make([]string, 0)
		filePaths = append(filePaths, tFPath)
		filePaths = append(filePaths, layoutsList...)
		if core.Conf.Debug {
			log.Println(filePaths)
		}
		tempNew := templateCache{}
		tempNew.Expire = time.Now().Unix() + core.Conf.App.Web.TemplateTimeout
		var err error
		tempNew.Template, err = template.New(tmpName).Funcs(funcMap).ParseFiles(filePaths...)
		if err != nil {
			log.Printf("*Fail to load template:%s\n", tFPath)
			log.Println(err)
			return nil
		}

		cache[tmpName] = tempNew
		temp = tempNew
	}
	gettext.SetLocale(locale)
	return temp.Template
}

func initI18n(locale string, domain string, dir string) {
	if core.Conf.Debug {
		log.Println("Init lang with", locale, domain, dir)
	}
	gettext.SetLocale(locale)
	gettext.Textdomain(domain)
	gettext.BindTextdomain(domain, dir, nil)
}

func InitTemplate() {
	for _, locale := range core.Conf.App.Web.Locales {
		initI18n(locale, "ui", core.Conf.App.Web.I18nPath)
	}
	// Load Layout List;
	lDir, err := ioutil.ReadDir(core.Conf.App.Web.LayoutsPath)
	if err == nil {
		for _, lFile := range lDir {
			if !lFile.IsDir() && filepath.Ext(lFile.Name()) == ".html" {
				lPath := filepath.Join(core.Conf.App.Web.LayoutsPath, lFile.Name())
				layoutsList = append(layoutsList, lPath)
			}
		}
	}
	cache = make(map[string]templateCache)
}
