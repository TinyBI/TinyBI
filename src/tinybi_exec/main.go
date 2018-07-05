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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"os"
	"os/signal"
	"plugin"
	"runtime"
	"syscall"
	"tinybi/core"
	"tinybi/mailer"
	"tinybi/model"
	"tinybi/task"
	"tinybi/web"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func main() {
	configPath := flag.String("c", "etc/config.json", "Path of configuration file")
	logPath := flag.String("l", "stdout", "Path of log file, use stdout to print logs at console")
	showHelp := flag.Bool("h", false, "Show usage")
	flag.Parse()
	if *showHelp {
		flag.Usage()
		return
	}
	//Init Log;
	var fLog *os.File
	if *logPath != "stdout" {
		//Redirect log to file;
		fLog, err := os.OpenFile(*logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Println("Fail to open log file", *logPath)
			log.Fatal(err)
		}
		log.SetOutput(fLog)
	}
	log.Println("ECBI Report System, version: 3.0")
	log.Println("Development Code: May Flowers")
	log.Println("*Starting ECBI daemon......................")
	//Init signals;
	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func(fToClose *os.File) {
		sigCatch := <-sChan
		log.Println("Catched signal:", sigCatch)
		if fToClose != nil {
			fToClose.Close()
		}
		log.Println("*ECBI daemon is stopped")
		os.Exit(0)
	}(fLog)
	initConf(*configPath)
	initDB()
	initScheduler()
	initWebMods()
	startScheduler()
	startMailQueue()
	startWebServer()
}

func initConf(confPath string) {
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
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs + 1)
}

func initDB() {
	var err error
	core.DB, err = xorm.NewEngine(core.Conf.DB.Driver, core.Conf.DB.Connection)
	if err != nil {
		log.Println("Fail to connect to Database")
		log.Fatal(err)
	}
	err = core.DB.Ping()
	if err != nil {
		log.Println("Fail to connect to Database")
		log.Fatal(err)
	}
	initDBTables()
}

func initDBTables() {
	model.InstallCoreTaskTable()
	model.InstallCoreSettingsTable()
	model.InstallAclTables()
	model.InstallModPhpTables()
	//model.InstallTmallTagsTables()
}

func initScheduler() {
	task.RegScheduleTasks()
}

func startScheduler() {
	task.StartSchedule()
}

func startMailQueue() {
	if core.Conf.MailQueue {
		mailer.Start()
	}
}

func startWebServer() {
	web.InitTemplate()
	http.HandleFunc("/", HttpServer)
	log.Fatal(http.ListenAndServe(core.Conf.Web.Host, nil))
}

func HttpServer(w http.ResponseWriter, r *http.Request) {
	if core.Conf.Debug {
		log.Println(r.URL.Path)
	}
	publicPrefix := fmt.Sprintf("/%s/", core.ConfPublicFolder)
	publicFilePath := filepath.Join(core.Conf.Web.RootPath, core.ConfPublicFolder, strings.Replace(r.URL.Path, publicPrefix, "", 1))
	if strings.HasPrefix(r.URL.Path, publicPrefix) && !strings.HasSuffix(r.URL.Path, "/") {
		http.ServeFile(w, r, publicFilePath)
	} else {
		routePath := r.URL.Path
		if routePath == "/" {
			routePath = "/index.html"
		}
		app, ok := web.WebRoutes[routePath]
		if !ok {
			web.ErrorNotFound(w, r)
		} else {
			app.Dispatch(w, r)
		}
	}
}

func initWebMods() {
	mDir, err := ioutil.ReadDir(core.Conf.ModsPath)
	if err == nil {
		for _, mFile := range mDir {
			if !mFile.IsDir() && filepath.Ext(mFile.Name()) == ".so" {
				mPath := filepath.Join(core.Conf.ModsPath, mFile.Name())
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
						webApp, ok := plugWebApp.(web.WebApp)
						if !ok {
							if core.Conf.Debug {
								log.Println("The mod does not implement WebApp:", mPath)
							}
						} else {
							route := "/"
							route += strings.TrimSuffix(mFile.Name(), filepath.Ext(mFile.Name()))
							route += ".html"
							web.WebRoutes[route] = webApp
							if core.Conf.Debug {
								log.Println("Loaded module with route", route, "from", mFile.Name())
							}
							log.Printf("Loaded module:%s, version:%s\n", webApp.Version().Name, webApp.Version().Version)
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
