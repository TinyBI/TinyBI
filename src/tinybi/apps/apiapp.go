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
package apps

import (
	"net/http"
	"strings"
	"tinybi/restful"
	"tinybi/webcore"
)

type ApiApp struct {
	webcore.BaseWebApp
}

func (this ApiApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	if len(paths) < 2 {
		webcore.ErrorNotFound(w, r)
		return
	}
	resourceName := paths[1]
	resource, ok := restful.RestfulResources[resourceName]
	if !ok {
		webcore.ErrorNotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		resource.Get(w, r)
		break
	case http.MethodPost:
		resource.Post(w, r)
		break
	case http.MethodPut:
		resource.Put(w, r)
		break
	case http.MethodDelete:
		resource.Delete(w, r)
		break
	default:
		webcore.ErrorNotFound(w, r)
		return
	}
}
