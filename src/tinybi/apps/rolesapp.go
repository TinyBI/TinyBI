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
	"encoding/json"
	"log"
	"net/http"
	"tinybi/core"
	"tinybi/webcore"
)

type RolesApp struct {
	webcore.BaseWebApp
}

type roleRow struct {
	RoleId    int64  `json:"0"`
	RoleName  string `json:"1"`
	RoleCodes string `json:"2"`
}

func (this RolesApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		webcore.AclCheckRedirect(w, r, "SYSTEM", "/login.html")
		switch r.URL.Query().Get("act") {
		case "list":
			//List concurrent tasks which are owned by current user;
			//Super Administrators can see all tasks;
			this.list(w, r)
			break
		default:
			//Show WEB page;
			this.showPage(w, r)
			break
		}
	} else {
		webcore.AclCheckRedirect(w, r, "SYSTEM", "/login.html")
		switch r.URL.Query().Get("act") {
		case "addRole":
			this.addUser(w, r)
			break
		case "delRole":
			this.delUser(w, r)
			break
		}
	}
}

func (this RolesApp) showPage(w http.ResponseWriter, r *http.Request) {
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "role_manager.html").Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func (this RolesApp) list(w http.ResponseWriter, r *http.Request) {
	//AJAX Method;
	nullRet := `{"data":[]}`
	var fullRet struct {
		Data []roleRow `json:"data"`
	}
	//Return JSON Data;
	sql := "SELECT id, role_name, role_codes FROM core_roles"

	row, err := core.DB.Query(sql)
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query SQL", err)
		}
		w.Write([]byte(nullRet))
		return
	}
	defer row.Close()
	urs := make([]roleRow, 0)
	for row.Next() {
		ur := roleRow{}
		err = row.Scan(&ur.RoleId, &ur.RoleName, &ur.RoleCodes)
		if err != nil {
			w.Write([]byte(nullRet))
			return
		}
		urs = append(urs, ur)
	}
	fullRet.Data = urs
	sret, err := json.Marshal(fullRet)
	if err != nil {
		w.Write([]byte(nullRet))
		return
	}
	w.Write(sret)
}

func (this RolesApp) addUser(w http.ResponseWriter, r *http.Request) {
	//
}

func (this RolesApp) delUser(w http.ResponseWriter, r *http.Request) {
	//
}
