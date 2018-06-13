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
package task

import (
	"bytes"
	"os/exec"
	"sync"
	"tinybi/logger"
	"tinybi/model"
)

//Run PHP mod;
//Compatible with old php modules;
//Since all setters return the pointer of handler, a typical task registeration would be like this:
//RegTasks["TASK_NAME"]=newModPhpHandler().SetFilePath(path).SetParameters(params)
type modPhpHandler struct {
	loggerId   string
	phpExec    string
	filePath   string
	parameters []string
	mailScript string
	BaseHandler
}

func newModPhpHandler() *modPhpHandler {
	handler := new(modPhpHandler)
	handler.Mutex = new(sync.Mutex)
	handler.loggerId = "MOD_PHP_EXEC"
	phpExec := model.BusinessSettings.Get("PHP_EXEC")
	if phpExec != nil {
		handler.phpExec = phpExec.Value
	}
	if handler.phpExec == "" {
		handler.phpExec = "/usr/local/php/bin/php"
	}
	return handler
}

func (this *modPhpHandler) SetFilePath(path string) *modPhpHandler {
	this.filePath = path
	return this
}

func (this *modPhpHandler) SetParameters(params ...string) *modPhpHandler {
	this.parameters = make([]string, 0)
	this.parameters = append(this.parameters, params...)
	return this
}

func (this *modPhpHandler) SetMailScript(path string) *modPhpHandler {
	this.mailScript = path
	return this
}

func (this modPhpHandler) Exec() {
	defer func() {
		if err := recover(); err != nil {
			logger.Printf(this.loggerId, "*Panic:%s", err)
		}
	}()
	if this.filePath == "" {
		return
	}
	logger.Printf(this.loggerId, "Start of php module:%s", this.filePath)
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	args := make([]string, 0)
	args = append(args, this.filePath)
	if len(this.parameters) > 0 {
		args = append(args, this.parameters...)
	}
	cmd := exec.Command(this.phpExec, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Printf(this.loggerId,
			"*Error occurs when running (%s), php script:%s, parameters:%s",
			err.Error(), this.filePath, this.parameters)
	} else {
		if this.mailScript != "" {
			mailCmd := exec.Command("/bin/sh", this.mailScript)
			mailCmd.Stdout = &out
			err = mailCmd.Run()
			if err != nil {
				logger.Printf(this.loggerId,
					"*Error occurs sending mail after run php module (%s), mail script:%s",
					err.Error(), this.mailScript)
			}
		}
	}
	logger.Printf(this.loggerId, out.String())
	logger.Printf(this.loggerId, "End of php module:%s", this.filePath)
}
