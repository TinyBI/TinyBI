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
package models

import (
	"errors"
	"log"
	"net/http"
	"time"
	"tinybi/core"
	"tinybi/webcore"
)

//Current Tasks;
type ConcurrentTask struct {
	Id          int64     `xorm:"'id'"`
	Description string    `xorm:"'description'"`
	Status      string    `xorm:"'status'"`
	Percentage  float32   `xorm:"'percentage'"`
	StartTime   string    `xorm:"'start_time'"`
	EndTime     string    `xorm:"'end_time'"`
	LastUpdated time.Time `xorm:"'last_updated'"`
	Owner       string    `xorm:"'owner'"`
	OwnerId     int64     `xorm:"'owner_id'"`
	Error       string    `xorm:"'error'"`
	FilePath    string    `xorm:"'file_path'"`
}

type TasksModel struct {
	//Operation Model;
}

var WebTaskModel *TasksModel

func init() {
	WebTaskModel = new(TasksModel)
}

const TaskStatusPending = "PENDING"
const TaskStatusRunning = "RUNNING"
const TaskStatusWarning = "WARNING"
const TaskStatusError = "ERROR"
const TaskStatusDone = "DONE"

//New User Level Concurrent Task;
func (this *TasksModel) NewTask(r *http.Request) *ConcurrentTask {
	session := webcore.AclGetSession(r)
	if session == nil {
		return nil
	}
	task := new(ConcurrentTask)
	task.Status = TaskStatusPending
	task.OwnerId = session.User.Id
	task.Owner = session.User.NickName
	return task
}

//New System Level Concurrent Task;
func (this *TasksModel) NewSystemTask() *ConcurrentTask {
	task := new(ConcurrentTask)
	task.Status = TaskStatusPending
	task.OwnerId = 0
	task.Owner = "System"
	return task
}

func (this *ConcurrentTask) Push() error {
	if this.Id != 0 {
		return errors.New("The task has already been pushed")
	}
	insertId, err := core.DBEngine.Table("core_concurrent_tasks").Insert(this)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
		}
	}
	this.Id = insertId
	return nil
}

func (this *ConcurrentTask) SetPercentage(p float32) error {
	this.Percentage = p
	return this.Save()
}

func (this *ConcurrentTask) Start() error {
	this.StartTime = time.Now().Format("2006-01-02 15:04:05")
	this.Status = TaskStatusRunning
	return this.Save()
}

func (this *ConcurrentTask) Done() error {
	this.EndTime = time.Now().Format("2006-01-02 15:04:05")
	this.Status = TaskStatusDone
	return this.Save()
}

func (this *ConcurrentTask) Save() error {
	if this.Id == 0 {
		return errors.New("The task has not been pushed")
	}
	this.LastUpdated = time.Now()
	_, err := core.DBEngine.Table("core_concurrent_tasks").Where("id=?", this.Id).Update(this)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
		}
	}
	return nil
}
