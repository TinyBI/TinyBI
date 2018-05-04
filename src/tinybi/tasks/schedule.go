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
	"strconv"
	"sync"
	"tinybi/core"
	"tinybi/models"
)

type Handler interface {
	Init()
	Exec(...interface{})
}

type BaseHandler struct {
	mutex *sync.Mutex
}

func (this BaseHandler) Init() {
	this.mutex = new(sync.Mutex)
}

func (this BaseHandler) Exec(...interface{}) {
	//
}

//Defined tasks below
//Steps to run a task:
//I, create a object (tasker) that implements tasks.Handler;
//II, register it with tasks.RegTasks[name]=tasker (at reg.go);
//III, write a record to core_tasks;

var RegTasks map[string]Handler

func init() {
	RegTasks = make(map[string]Handler)
}

func ReloadScheduledTasks() {
	//Unload all tasks first;
	core.Scheduler.Clear()
	//Read configurations from DB;
	var tasks []models.CoreTasks
	err := core.DBEngine.Table("core_tasks").Where("enabled='YES'").Find(&tasks)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
			return
		}
	}
	for _, task := range tasks {
		if core.Conf.Debug {
			log.Println(task)
		}
		rTask, ok := RegTasks[task.TaskName]
		if ok {
			switch task.ScheduleType {
			case "SECONDS":
				interval, err := strconv.Atoi(task.ScheduleAt)
				if err == nil {
					core.Scheduler.Every(uint64(interval)).Seconds().Do(func() { rTask.Exec() })
				}
				if core.Conf.Debug {
					log.Println("Installed task", task.TaskName, "for", interval, "seconds")
				}
				break
			case "MINUTES":
				interval, err := strconv.Atoi(task.ScheduleAt)
				if err == nil {
					core.Scheduler.Every(uint64(interval)).Minutes().Do(func() { rTask.Exec() })
				}
				if core.Conf.Debug {
					log.Println("Installed task", task.TaskName, "for", interval, "minutes")
				}
				break
			case "HOURLY":
				core.Scheduler.Every(1).Hour().At(task.ScheduleAt).Do(func() { rTask.Exec() })
				break
			case "DAILY":
				core.Scheduler.Every(1).Day().At(task.ScheduleAt).Do(func() { rTask.Exec() })
				break
			case "MONDAY":
				core.Scheduler.Every(1).Monday().At(task.ScheduleAt).Do(func() { rTask.Exec() })
				break
			case "TUESDAY":
				core.Scheduler.Every(1).Tuesday().At(task.ScheduleAt).Do(func() { rTask.Exec() })
				break
			case "WEDNESDAY":
				core.Scheduler.Every(1).Wednesday().At(task.ScheduleAt).Do(func() { rTask.Exec() })
				break
			case "THURSDAY":
				core.Scheduler.Every(1).Thursday().At(task.ScheduleAt).Do(func() { rTask.Exec() })
				break
			case "FRIDAY":
				core.Scheduler.Every(1).Friday().At(task.ScheduleAt).Do(func() { rTask.Exec() })
				break
			case "SATURDAY":
				core.Scheduler.Every(1).Saturday().At(task.ScheduleAt).Do(func() { rTask.Exec() })
				break
			case "SUNDAY":
				core.Scheduler.Every(1).Sunday().At(task.ScheduleAt).Do(func() { rTask.Exec() })
				break
			}
		}
	}
}
