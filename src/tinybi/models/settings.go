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
	"time"
	"tinybi/core"
)

//Core Business Settings;
type Settings struct {
	Id          int64
	Code        string    `xorm:"'code'" json:"1"`
	Description string    `xorm:"'description'" json:"2"`
	Value       string    `xorm:"'value'" json:"3"`
	LastUpdated time.Time `xorm:"'last_updated' default 'CURRENT_tIMESTAMP'" json:"4"`
}

type SettingsModel struct {
	//Operation Model;
	//Using cache to decrease I/O;
	cache map[string]*Settings
}

var BusinessSettings *SettingsModel

func init() {
	BusinessSettings = new(SettingsModel)
	BusinessSettings.Init()
}

func (this *SettingsModel) Init() {
	this.cache = make(map[string]*Settings)
}

func (this *SettingsModel) Get(code string) *Settings {
	var settings *Settings
	settings, ok := this.cache[code]
	if !ok {
		//Try to get the data from DB;
		_, err := core.DBEngine.Table("core_settings").Where("code=?", code).Get(settings)
		if err != nil && core.Conf.Debug {
			log.Println(err)
			return nil
		}
	}
	return settings
}

func (this *SettingsModel) Set(settings *Settings) error {
	if settings == nil {
		return errors.New("Null pointer")
	}
	if settings.Id == 0 {
		insertId, err := core.DBEngine.Table("core_settings").Insert(settings)
		if err != nil {
			if core.Conf.Debug {
				log.Println(err)
			}
			return err
		}
		settings.Id = insertId
	} else {
		affects, err := core.DBEngine.Table("core_settings").Where("id=?", settings.Id).Update(settings)
		if err != nil {
			if core.Conf.Debug {
				log.Println(err)
			}
			return err
		}
		if affects == 0 {
			return errors.New("No change for the settings")
		}
	}
	this.cache[settings.Code] = settings
	return nil
}

func (this *SettingsModel) List() []Settings {
	list := make([]Settings, 0)
	for _, s := range this.cache {
		setting := Settings{Id: s.Id, Code: s.Code,
			Description: s.Description, Value: s.Value,
			LastUpdated: s.LastUpdated}
		list = append(list, setting)
	}
	return list
}
