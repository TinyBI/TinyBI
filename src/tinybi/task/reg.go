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
	"log"
	"strings"
	"tinybi/core"
	"tinybi/model"
)

//Register tasks here;
func RegScheduleTasks() {
	//Load php modules;
	regModPhpTasks()
}

//Register php mods;
func regModPhpTasks() {
	var modPhps []model.ModPhp
	err := core.DB.Table("core_mod_phps").Where("1=1").Find(&modPhps)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
		}
		return
	}
	for _, modPhp := range modPhps {
		taskName := "MOD_PHP_"
		taskName += modPhp.Code
		phpHandler := newModPhpHandler().SetFilePath(modPhp.FilePath)
		if modPhp.Parameters != "" {
			parameters := strings.Split(modPhp.Parameters, ",")
			phpHandler.SetParameters(parameters...)
		}
		if modPhp.MailScript != "" {
			phpHandler.SetMailScript(modPhp.MailScript)
		}
		RegTasks[taskName] = phpHandler
	}
}
