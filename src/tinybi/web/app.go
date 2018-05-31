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
package web

import "net/http"

type WebApp interface {
	Version() VersionInfo
	Dispatch(w http.ResponseWriter, r *http.Request)
}

type VersionInfo struct {
	Name        string
	Version     string
	Description string
}

type BaseWebApp struct {
	AppVersion *VersionInfo
}

func NewBaseApp() *BaseWebApp {
	app := new(BaseWebApp)
	app.AppVersion = new(VersionInfo)
	return app
}

func (this BaseWebApp) Version() VersionInfo {
	return *this.AppVersion
}

func (this BaseWebApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

var WebRoutes map[string]WebApp

func init() {
	WebRoutes = make(map[string]WebApp)
}
