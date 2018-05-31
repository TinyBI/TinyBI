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
package core

const ConfPublicFolder string = "public"
const ConfTemplatesFolder string = "templates"
const ConfLayoutsFolder string = "layouts"
const ConfI18nFolder string = "i18n"
const ConfMenusFolder string = "menus"

type DbConf struct {
	Driver     string `json:"driver"`
	Connection string `json:"connection"`
}

type WebConf struct {
	Host            string   `json:"host"`
	SessionTimeout  int64    `json:"session_timeout"`
	TemplateTimeout int64    `json:"template_timeout"`
	Locales         []string `json:"locales"`
	RootPath        string   `json:"root_path"`
}

type LogConf struct {
	RootPath string `json:"root_path"`
}

type Configuration struct {
	DB        DbConf  `json:"db"`
	Web       WebConf `json:"web"`
	Logger    LogConf `json:"logger"`
	TimeZone  string  `json:"time_zone"`
	Debug     bool    `json:"debug"`
	MailQueue bool    `json:"mail_queue"`
	ModsPath  string  `json:"mods_path"`
}

var Conf Configuration

func init() {
	Conf = Configuration{}
}
