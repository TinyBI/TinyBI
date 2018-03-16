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
package webcore

import (
	"net/http"
	"strings"
)

// All WEB App MUST implement the interface below;

type RestfulResource interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type BaseResource struct {
	//
}

func (this BaseResource) Get(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (this BaseResource) Post(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (this BaseResource) Put(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (this BaseResource) Delete(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

//Parsing URI into parameters;
//e.g. /API/Node/2
//Return: params["Node"] => 2;
func (this BaseResource) ParseParams(r *http.Request) map[string]string {
	rawParams := strings.Split(r.URL.Path, "/")
	params := make(map[string]string)
	for i := 1; i < len(rawParams); i += 2 {
		if i+1 < len(rawParams) {
			params[rawParams[i]] = rawParams[i+1]
		} else {
			params[rawParams[i]] = ""
		}
	}
	return params
}
