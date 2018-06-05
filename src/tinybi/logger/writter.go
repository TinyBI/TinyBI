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
package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tinybi/core"
)

func Printf(id string, format string, data ...interface{}) error {
	pathPrefix := core.Conf.Logger.RootPath
	if pathPrefix == "" {
		return errors.New("You must set the root path of logger to enable it")
	}
	curFolder := core.FromUnixTime(time.Now().Unix(), core.DefaultDateFormat)
	path := filepath.Join(pathPrefix, curFolder)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	path = filepath.Join(pathPrefix, curFolder, id)
	path += ".txt"
	logStr := format
	if len(data) > 0 {
		logStr = fmt.Sprintf(format, data)
	}
	writeStr := fmt.Sprintf("[%s]%s\n", core.NowTime(), logStr)
	fLog, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fLog.Close()
	_, err = fLog.WriteString(writeStr)
	return err
}

func GetPath(id string) string {
	pathPrefix := core.Conf.Logger.RootPath
	if pathPrefix == "" {
		return ""
	}
	curFolder := core.FromUnixTime(time.Now().Unix(), core.DefaultDateFormat)
	path := filepath.Join(pathPrefix, curFolder)
	path = filepath.Join(pathPrefix, curFolder, id)
	path += ".txt"
	return path
}
