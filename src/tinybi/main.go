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
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"plugin"
	"strings"
	"tinybi/apps"
	"tinybi/core"
	"tinybi/models"
	"tinybi/tasks"
	"tinybi/webcore"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func main() {
	configPath := flag.String("c", "etc/config.json", "Path of configuration file")
	flag.Parse()
	initApp(*configPath)
	if core.Conf.Debug {
		log.Println(core.Conf)
	}
	initData()
	initScheduler()
	initMods()
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
	//Use ORM instead of native DB Engine;
	core.DBEngine, err = xorm.NewEngine(core.Conf.DB.Driver, core.Conf.DB.Connection)
	if err != nil {
		log.Println("Fail to connect to Database")
		log.Fatal(err)
	}
	//core.DB.SetMaxIdleConns(100)
	err = core.DBEngine.Ping()
	if err != nil {
		log.Println("Fail to connect to Database")
		log.Fatal(err)
	}
	core.DB = core.DBEngine.DB().DB
}

func initScheduler() {
	tasks.RegScheduleTasks()
	tasks.StartSchedule()
}

func initData() {
	//Init master data for modules;
	//General Ledger;
	glModel := models.GLModel{}
	glAccountPath := filepath.Join(core.Conf.Data.MasterPath, "gl_accounts.json")
	err := glModel.InitMasterAccounts(glAccountPath)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
		}
	}
	glModel.InitMasterPeriods()
}

func initMods() {
	mDir, err := ioutil.ReadDir(core.Conf.App.Web.ModsPath)
	if err == nil {
		for _, mFile := range mDir {
			if !mFile.IsDir() && filepath.Ext(mFile.Name()) == ".so" {
				mPath := filepath.Join(core.Conf.App.Web.ModsPath, mFile.Name())
				if core.Conf.Debug {
					log.Println("Load module from:", mPath)
				}
				plug, err := plugin.Open(mPath)
				if err == nil {
					plugWebApp, err := plug.Lookup("ModWebApp")
					if err != nil {
						if core.Conf.Debug {
							log.Println("Fail to load module from:", mPath)
							log.Println(err)
						}
					} else {
						webApp, ok := plugWebApp.(webcore.WebApp)
						if !ok {
							if core.Conf.Debug {
								log.Println("Fail to load module from:", mPath)
							}
						} else {
							route := "/mod_"
							route += strings.TrimSuffix(mFile.Name(), filepath.Ext(mFile.Name()))
							route += ".html"
							apps.WebRoutes[route] = webApp
							if core.Conf.Debug {
								log.Println("Loaded module with route", route, "from", mFile.Name())
							}
						}
					}
				} else {
					if core.Conf.Debug {
						log.Println("Fail to load module from:", mPath)
						log.Println(err)
					}
				}
			}
		}
	}
}
