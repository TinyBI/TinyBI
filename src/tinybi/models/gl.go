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
	"time"
)

//General Ledgers;
//Accounting Periods;
type GLPeriod struct {
	Id          int64     `xorm:"'id'"`
	PeriodCode  string    `xorm:"'period_code'"`
	PeriodName  string    `xorm:"'period_name'"`
	Status      string    `xorm:"'status' default 'CLOSED'"`
	Description string    `xorm:"'description'"`
	StartTime   int       `xorm:"'start_time'"`
	EndTime     int       `xorm:"'end_time'"`
	LastUpdated time.Time `xorm:"'last_updated' default 'CURRENT_tIMESTAMP'"`
}

type GLModel struct {
	//Operation Model;
	//Business Operations;
}
