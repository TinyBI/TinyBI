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
	"errors"
	"net/http"
	"sync"
	"time"
	"tinybi/core"
)

type TaskRunMethod func(*Task) error

type Task struct {
	Id          int64
	Description string
	Status      string
	Percentage  float32
	StartTime   string
	EndTime     string
	Owner       string
	OwnerId     int64
	FilePath    string
	Error       string
	Run         TaskRunMethod
}

type TaskMaster struct {
	Mutex     *sync.RWMutex
	Tasks     map[int64]*Task
	UserTasks map[int64]int //Number of active tasks per user has;
}

const TaskStatusPending = "PENDING"
const TaskStatusRunning = "RUNNING"
const TaskStatusWarning = "WARNING"
const TaskStatusError = "ERROR"
const TaskStatusDone = "DONE"

var WebTaskMaster *TaskMaster

func init() {
	WebTaskMaster = new(TaskMaster)
	WebTaskMaster.Tasks = make(map[int64]*Task)
	WebTaskMaster.Mutex = new(sync.RWMutex)
	WebTaskMaster.UserTasks = make(map[int64]int)
}

//New User Level Concurrent Task;
func NewTask(r *http.Request) *Task {
	session := AclGetSession(r)
	if session == nil {
		return nil
	}
	task := new(Task)
	task.Status = TaskStatusPending
	task.OwnerId = session.User.Id
	task.Owner = session.User.NickName
	return task
}

//New System Level Concurrent Task;
func NewSystemTask() *Task {
	task := new(Task)
	task.Status = TaskStatusPending
	task.OwnerId = 0
	task.Owner = "System"
	return task
}

//Push Task into DB & Master;
func (this *Task) Push() error {
	if this.Status != TaskStatusPending {
		return errors.New("Only pending tasks can be pushed")
	}
	if WebTaskMaster.OwnerTasks(this.OwnerId) >= core.Conf.App.Web.MaxTasksPerUser {
		return errors.New("Reach the max number of active task per user")
	}
	sql := "INSERT INTO core_concurrent_tasks "
	sql += "(description,status,percentage,owner,owner_id) "
	sql += "(?,?,0,?,?) "
	res, err := core.DB.Exec(sql, this.Description, this.Status, this.Owner, this.OwnerId)
	if err != nil {
		return err
	}
	this.Id, err = res.LastInsertId()
	if err != nil {
		return err
	}
	WebTaskMaster.Push(this)
	return nil
}

//Remove Task from Master;
func (this *Task) Pull() {
	if this.Id != 0 {
		WebTaskMaster.Pull(this)
	}
}

//Set Running Pencentage;
func (this *Task) SetPercentage(p float32) error {
	if this.Id == 0 {
		return errors.New("The task has not been pushed")
	}
	this.Percentage = p
	sql := "UPDATE core_concurrent_tasks "
	sql += "SET percentage = ? "
	sql += "WHERE id=?"
	_, err := core.DB.Exec(sql, this.Percentage, this.Id)
	if err != nil {
		return err
	}
	return nil
}

//Set Status;
//Parameters:
//status
//filePath
//error
func (this *Task) SetStatus(s ...string) error {
	if this.Id == 0 {
		return errors.New("The task has not been pushed")
	}
	if len(s) == 0 {
		return errors.New("Illegal call of Task::SetStatus")
	}
	status := s[0]
	if status == TaskStatusError && len(s) < 3 {
		return errors.New("An error task must has both status, output path(can be empty) and error description")
	}
	this.Status = status
	sql := "UPDATE core_concurrent_tasks "
	sql += "SET status = ? "
	switch status {
	case TaskStatusRunning:
		sql += ",stat_time = now() "
		sql += ",percentage = 0 "
		break
	case TaskStatusError:
		//Error task must set its error message;
		sql += ",error = "
		sql += s[2]
		sql += " "
	case TaskStatusWarning:
	case TaskStatusDone:
		sql += ",end_time = now() "
		sql += ",percentage = 100 "
		//Set output path;
		if len(s) >= 2 && s[1] != "" {
			sql += ",file_path = '"
			sql += s[1]
			sql += "' "
		}
		this.Pull()
		break
	}
	sql += "WHERE id=?"
	_, err := core.DB.Exec(sql, this.Status, this.Id)
	if err != nil {
		return err
	}
	return nil
}

//Brief call when task is done;
func (this *Task) Done() error {
	return this.SetStatus(TaskStatusDone)
}

//Run the task;
func (this *Task) SyncRun(t *Task) error {
	if this.Status != TaskStatusPending {
		return errors.New("Only pending tasks can be run")
	}
	this.SetStatus(TaskStatusRunning)
	err := this.Run(t)
	if err != nil {
		this.Error = err.Error()
		return err
	}
	return nil
}

func (this *Task) AsyncRun(t *Task) {
	go this.Run(t)
}

//Task Schedule;
func (this *TaskMaster) Run() {
	for {
		if len(this.Tasks) > 0 {
			for _, t := range this.Tasks {
				if t.Status == TaskStatusPending {
					t.AsyncRun(t)
				}
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func (this *TaskMaster) Start() {
	go this.Run()
}

func (this *TaskMaster) Push(t *Task) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.Tasks[t.Id] = t
	nT, ok := this.UserTasks[t.OwnerId]
	if !ok {
		nT = 0
	}
	nT += 1
	this.UserTasks[t.OwnerId] = nT
}

func (this *TaskMaster) Pull(t *Task) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	delete(this.Tasks, t.Id)
	nT, ok := this.UserTasks[t.OwnerId]
	if !ok {
		nT = 1
	}
	nT -= 1
	this.UserTasks[t.OwnerId] = nT
}

func (this *TaskMaster) OwnerTasks(onwerId int64) int {
	if onwerId == 0 {
		return 9999999
	}
	nT, ok := this.UserTasks[onwerId]
	if !ok {
		nT = 0
	}
	return nT
}
