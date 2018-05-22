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
package tasks

import (
	"log"
	"reflect"
	"strconv"
	"sync"

	"tinybi/core"
	"tinybi/models"
)

const ReloadInterval uint64 = 10

type handler interface {
	IsTaskUpdated(models.CoreTasks) bool
	IsScheduled() bool
	Exec()
	GetExecAddr() func()
}

type baseHandler struct {
	execAddr  func()
	task      models.CoreTasks
	scheduled bool
	mutex     *sync.Mutex
}

func newHandler() *baseHandler {
	handler := new(baseHandler)
	handler.mutex = new(sync.Mutex)
	return handler
}

func (this baseHandler) Exec() {
	//
}

func (this baseHandler) GetExecAddr() func() {
	return this.execAddr
}

func (this *baseHandler) SetExecAddr(f func()) {
	this.execAddr = f
}

func (this baseHandler) IsScheduled() bool {
	return this.scheduled
}

func (this baseHandler) IsTaskUpdated(task models.CoreTasks) bool {
	if this.task.LastUpdated.Unix() < task.LastUpdated.Unix() ||
		this.task.Id != task.Id {
		return true
	}
	return false
}

func (this *baseHandler) SetScheduled(isScheduled bool) {
	this.scheduled = isScheduled
}

func (this *baseHandler) UpdateTask(task models.CoreTasks) {
	this.task.Id = task.Id
	this.task.TaskName = task.TaskName
	this.task.Description = task.Description
	this.task.Enabled = task.Enabled
	this.task.ScheduleType = task.ScheduleType
	this.task.ScheduleAt = task.ScheduleAt
	this.task.LastUpdated = task.LastUpdated
}

func SetExecAddr(this handler, f func()) {
	setMethod := reflect.ValueOf(this).MethodByName("SetExecAddr")
	if setMethod.IsValid() {
		params := make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(f)
		setMethod.Call(params)
	}
}

func SetScheduled(this handler, isScheduled bool) {
	setMethod := reflect.ValueOf(this).MethodByName("SetScheduled")
	if setMethod.IsValid() {
		params := make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(isScheduled)
		setMethod.Call(params)
	}
}

func UpdateTask(this handler, task models.CoreTasks) {
	setMethod := reflect.ValueOf(this).MethodByName("UpdateTask")
	if setMethod.IsValid() {
		params := make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(task)
		setMethod.Call(params)
	}
}

//Defined tasks below
//Steps to run a task:
//I, create a object (tasker) that implements tasks.handler;
//II, register it with tasks.RegTasks[name]=tasker (at reg.go);
//III, write a record to core_tasks;

var RegTasks map[string]handler

func init() {
	RegTasks = make(map[string]handler)
}

func ReloadScheduledTasks() {
	//Read configurations from DB;
	var tasks []models.CoreTasks
	err := core.DBEngine.Table("core_tasks").Where("1=1").Find(&tasks)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
			return
		}
	}
	for _, task := range tasks {
		rTask, ok := RegTasks[task.TaskName]
		if ok {
			if task.Enabled == "YES" {
				if rTask.IsScheduled() {
					if rTask.IsTaskUpdated(task) {
						if core.Conf.Debug {
							log.Println("Task", task.TaskName, "is updated:", task)
						}
						SetScheduled(rTask, false)
					}
				} else {
					log.Println("New unschedule task", task.TaskName, ":", task)
				}
			} else {
				//Disabled task;
				if rTask.IsScheduled() {
					log.Println("Task", task.TaskName, "has been disabled")
					core.Scheduler.Remove(rTask.GetExecAddr())
					SetScheduled(rTask, false)
				}
				continue
			}
			if !rTask.IsScheduled() {
				UpdateTask(rTask, task)
				SetExecAddr(rTask, func() {
					rTask.Exec()
				})
				switch task.ScheduleType {
				case "SECONDS":
					interval, err := strconv.Atoi(task.ScheduleAt)
					if err == nil {
						core.Scheduler.Every(uint64(interval)).Seconds().Do(rTask.GetExecAddr())
					}
					if core.Conf.Debug {
						log.Println("Installed task", task.TaskName, "for", interval, "seconds")
					}
					SetScheduled(rTask, true)
					break
				case "MINUTES":
					interval, err := strconv.Atoi(task.ScheduleAt)
					if err == nil {
						core.Scheduler.Every(uint64(interval)).Minutes().Do(rTask.GetExecAddr())
					}
					if core.Conf.Debug {
						log.Println("Installed task", task.TaskName, "for", interval, "minutes")
					}
					SetScheduled(rTask, true)
					break
				case "HOURLY":
					core.Scheduler.Every(1).Hour().At(task.ScheduleAt).Do(rTask.GetExecAddr())
					if core.Conf.Debug {
						log.Println("Installed task", task.TaskName, "at", task.ScheduleAt, "one hour")
					}
					SetScheduled(rTask, true)
					break
				case "DAILY":
					core.Scheduler.Every(1).Day().At(task.ScheduleAt).Do(rTask.GetExecAddr())
					if core.Conf.Debug {
						log.Println("Installed task", task.TaskName, "at", task.ScheduleAt, "of the day")
					}
					SetScheduled(rTask, true)
					break
				case "MONDAY":
					core.Scheduler.Every(1).Monday().At(task.ScheduleAt).Do(rTask.GetExecAddr())
					SetScheduled(rTask, true)
					break
				case "TUESDAY":
					core.Scheduler.Every(1).Tuesday().At(task.ScheduleAt).Do(rTask.GetExecAddr())
					SetScheduled(rTask, true)
					break
				case "WEDNESDAY":
					core.Scheduler.Every(1).Wednesday().At(task.ScheduleAt).Do(rTask.GetExecAddr())
					SetScheduled(rTask, true)
					break
				case "THURSDAY":
					core.Scheduler.Every(1).Thursday().At(task.ScheduleAt).Do(rTask.GetExecAddr())
					SetScheduled(rTask, true)
					break
				case "FRIDAY":
					core.Scheduler.Every(1).Friday().At(task.ScheduleAt).Do(rTask.GetExecAddr())
					SetScheduled(rTask, true)
					break
				case "SATURDAY":
					core.Scheduler.Every(1).Saturday().At(task.ScheduleAt).Do(rTask.GetExecAddr())
					SetScheduled(rTask, true)
					break
				case "SUNDAY":
					core.Scheduler.Every(1).Sunday().At(task.ScheduleAt).Do(rTask.GetExecAddr())
					SetScheduled(rTask, true)
					break
				}
			}
		} else {
			if core.Conf.Debug {
				log.Println("Unregister task:", task.TaskName)
			}
		}
	}
}

func StartSchedule() {
	core.Scheduler.Every(ReloadInterval).Seconds().Do(ReloadScheduledTasks)
	core.Scheduler.Start()
}
