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
	"database/sql"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"tinybi/apps"
	"tinybi/core"
	"tinybi/webcore"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configPath := flag.String("c", "etc/config.json", "Path of configuration file")
	flag.Parse()
	initApp(*configPath)
	if core.Conf.Debug {
		log.Println(core.Conf)
	}
	webcore.InitTemplate()
	http.HandleFunc("/", HttpServer)
	log.Fatal(http.ListenAndServe(core.Conf.App.Web.Host, nil))
}

func HttpServer(w http.ResponseWriter, r *http.Request) {
	if core.Conf.Debug {
		log.Println(r.URL.Path)
	}
	if strings.HasPrefix(r.URL.Path, "/public/") && !strings.HasSuffix(r.URL.Path, "/") {
		http.ServeFile(w, r, core.Conf.App.Web.PublicPath+strings.Replace(r.URL.Path, "/public", "", 1))
	} else {
		routePath := r.URL.Path
		app, ok := apps.WebRoutes[routePath]
		if !ok {
			http.Redirect(w, r, "/", http.StatusNotFound)
		} else {
			app.Dispatch(w, r)
		}
	}
}

func initApp(confPath string) {
	jsonStr, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Printf("Fail to load configuration from:%s\n", confPath)
		log.Fatal(err)
	}
	err = json.Unmarshal(jsonStr, &core.Conf)
	if err != nil {
		log.Printf("Fail to load configuration from:%s\n", confPath)
		log.Fatal(err)
	}
	core.DB, err = sql.Open(core.Conf.DB.Driver, core.Conf.DB.Connection)
	if err != nil {
		log.Println("Fail to connect to Database")
		log.Fatal(err)
	}
	//core.DB.SetMaxIdleConns(100)
	err = core.DB.Ping()
	if err != nil {
		log.Println("Fail to connect to Database")
		log.Fatal(err)
	}
}
