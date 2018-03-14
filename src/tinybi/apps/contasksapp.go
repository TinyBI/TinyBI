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
	"encoding/json"
	"log"
	"net/http"
	"tinybi/core"
	"tinybi/webcore"
)

type ConcurrentTasksApp struct {
	webcore.BaseWebApp
}

type concurrentTask struct {
	Id          int
	Description string
	Status      string
	Percentage  float32
	StartTime   string
	EndTime     string
	LastUpdated int64
	Owner       string
	OwnerId     int
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
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "concurrent_tasks.html").Execute(w, nil)
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
	//Return JSON Data;
	sql := "SELECT id, description, `status`, percentage, start_time, end_time, last_updated, `owner`, owner_id "
	sql += "FROM core_concurrent_tasks "
	sql += "WHERE owner_id = ?"
	session := webcore.AclGetSession(r)
	row, err := core.DB.Query(sql, session.User.Id)
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query SQL", err)
		}
		w.Write([]byte(nullRet))
		return
	}
	defer row.Close()
	tasks := make([]concurrentTask, 0)
	for row.Next() {
		task := concurrentTask{}
		err = row.Scan(&task.Id, &task.Description, &task.Status,
			&task.Percentage, &task.StartTime,
			&task.EndTime, &task.LastUpdated, &task.Owner,
			&task.OwnerId)
		if err != nil {
			w.Write([]byte(nullRet))
			return
		}
		tasks = append(tasks, task)
	}
	fullRet.Data = tasks
	sret, err := json.Marshal(fullRet)
	if err != nil {
		w.Write([]byte(nullRet))
		return
	}
	w.Write(sret)
}
