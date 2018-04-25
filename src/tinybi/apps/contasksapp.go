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
	"tinybi/models"
	"tinybi/webcore"
)

type ConcurrentTasksApp struct {
	webcore.BaseWebApp
}

type concurrentTask struct {
	Id          int64   `json:"0"`
	Description string  `json:"1"`
	Status      string  `json:"-"`
	Percentage  float32 `json:"2"`
	StartTime   string  `json:"3"`
	EndTime     string  `json:"4"`
	Owner       string  `json:"5"`
}

func (this ConcurrentTasksApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		webcore.AclCheckRedirect(w, r, "INDEX", "/login.html")
		switch r.URL.Query().Get("act") {
		case "list":
			//List concurrent tasks which are owned by current user;
			//Super Administrators can see all tasks;
			this.list(w, r)
			break
		default:
			//Show WEB page;
			this.showPage(w, r)
			break
		}
	} else {
		//This app does not accept POST action;
	}
}

func (this ConcurrentTasksApp) showPage(w http.ResponseWriter, r *http.Request) {
	var html struct {
		MaxTasksPerUser int
	}
	html.MaxTasksPerUser = core.Conf.App.Web.MaxTasksPerUser
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "concurrent_tasks.html").Execute(w, html)
	if err != nil {
		log.Println(err)
	}
}

func (this ConcurrentTasksApp) list(w http.ResponseWriter, r *http.Request) {
	//AJAX Method;
	nullRet := `{"data":[]}`
	var fullRet struct {
		Data []concurrentTask `json:"data"`
	}
	session := webcore.AclGetSession(r)
	if session == nil {
		w.Write([]byte(nullRet))
		return
	}
	var tasks []models.ConcurrentTask
	err := core.DBEngine.Table("core_concurrent_tasks").Where("owner_id=?", session.User.Id).Find(&tasks)
	if err != nil {
		w.Write([]byte(nullRet))
		return
	}
	for _, task := range tasks {
		var taskRow concurrentTask
		taskRow.Id = task.Id
		taskRow.Description = task.Description
		taskRow.EndTime = task.EndTime
		taskRow.Owner = task.Owner
		taskRow.Percentage = task.Percentage
		taskRow.StartTime = task.StartTime
		taskRow.Status = task.Status
		fullRet.Data = append(fullRet.Data, taskRow)
	}
	w.Write([]byte(webcore.JsonEncode(fullRet)))
}
