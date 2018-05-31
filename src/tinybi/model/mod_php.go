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
package model

import (
	"time"
	"tinybi/core"
)

type ModPhp struct {
	Id          int64
	Code        string    `xorm:"'code'"`
	Description string    `xorm:"'description'"`
	FilePath    string    `xorm:"'file_path' default ''"`
	Parameters  string    `xorm:" TEXT 'parameters' default ''"`
	MailScript  string    `xorm:"'mail_script' default ''"`
	LastUpdated time.Time `xorm:"TIMESTAMP 'last_updated' created updated" json:"4"`
}

func InstallModPhpTables() {
	core.DB.Table("core_mod_phps").Sync2(new(ModPhp))
}
