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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	Members   int    `json:"3"`
	EditLink  string `json:"4"`
	DelLink   string `json:"5"`
}

type roleDefine struct {
	Code        string `json:"code"`
	Description string `json:"description"`
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
		case "add":
			//Show "Add" page
			this.addPage(w, r)
			break
		case "edit":
			//Show "Edit" page
			this.editPage(w, r)
			break

		default:
			//Show WEB page;
			this.showPage(w, r)
			break
		}
	} else {
		webcore.AclCheckRedirect(w, r, "SYSTEM", "/login.html")
		switch r.URL.Query().Get("act") {
		case "add":
			this.add(w, r)
			break
		case "edit":
			this.edit(w, r)
			break
		case "del":
			this.del(w, r)
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
	sql := `SELECT
		id,
		role_name,
		role_codes,
		ifnull(a.members, 0) AS members
	FROM
		core_roles
	LEFT JOIN (
		SELECT
			role_id,
			count(1) AS members
		FROM
			core_users
		GROUP BY
			role_id
	) a ON core_roles.id = a.role_id`
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
		err = row.Scan(&ur.RoleId, &ur.RoleName, &ur.RoleCodes, &ur.Members)
		if err != nil {
			w.Write([]byte(nullRet))
			return
		}
		ur.EditLink = fmt.Sprintf("<p class='fa fa-edit'><a href='/roles.html?act=edit&id=%d'>Edit</a></p>", ur.RoleId)
		ur.DelLink = fmt.Sprintf("<p class='fa fa-trash-o'><a href='/roles.html?act=del&id=%d'>Delete</a></p>", ur.RoleId)
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

func (this RolesApp) addPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title       string
		Role        roleRow
		Act         string
		Acls        []roleDefine
		CheckedAcls []string
		Info        struct {
			Show    bool
			Message string
		}
	}
	//Load ACL definition from JSON file;
	jsonStr, err := ioutil.ReadFile(core.Conf.App.Web.AclDefinePath)
	if err != nil {
		log.Printf("Fail to load configuration from:%s\n", core.Conf.App.Web.AclDefinePath)
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	err = json.Unmarshal(jsonStr, &Html.Acls)
	if err != nil {
		log.Printf("Fail to load configuration from:%s\n", core.Conf.App.Web.AclDefinePath)
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	Html.Title = "New Role"
	Html.Act = "add"
	err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "role_manager_editor.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this RolesApp) add(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title       string
		Role        roleRow
		Act         string
		Acls        []roleDefine
		CheckedAcls []string
		Info        struct {
			Show    bool
			Message string
		}
	}
	r.ParseForm()
	Html.Role.RoleName = r.Form.Get("rolename")
	Html.Role.RoleCodes = strings.Join(r.Form["permissions"], ",")
	err := this.updateRole(&Html.Role)
	if err != nil {
		Html.Info.Show = true
		Html.Info.Message = "Faile to save the role"

		Html.Title = "New Role"
		Html.Act = "add"
		err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "role_manager_editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Redirect(w, r, "/roles.html", http.StatusFound)
	}
}

func (this RolesApp) editPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title       string
		Role        roleRow
		Act         string
		Acls        []roleDefine
		CheckedAcls []string
		Info        struct {
			Show    bool
			Message string
		}
	}
	//Load role info;
	roleId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of roles.html?act=edit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	sql := "SELECT id, role_name, role_codes FROM core_roles WHERE id=?"
	row, err := core.DB.Query(sql, roleId)
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query SQL", err)
		}
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&Html.Role.RoleId, &Html.Role.RoleName, &Html.Role.RoleCodes)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusNotFound)
			return
		}
	}
	Html.CheckedAcls = strings.Split(Html.Role.RoleCodes, ",")
	//Load ACL definition from JSON file;
	jsonStr, err := ioutil.ReadFile(core.Conf.App.Web.AclDefinePath)
	if err != nil {
		log.Printf("Fail to load configuration from:%s\n", core.Conf.App.Web.AclDefinePath)
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	err = json.Unmarshal(jsonStr, &Html.Acls)
	if err != nil {
		log.Printf("Fail to load configuration from:%s\n", core.Conf.App.Web.AclDefinePath)
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	Html.Title = "Edit Role"
	Html.Act = "edit"
	err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "role_manager_editor.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this RolesApp) edit(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title       string
		Role        roleRow
		Act         string
		Acls        []roleDefine
		CheckedAcls []string
		Info        struct {
			Show    bool
			Message string
		}
	}
	//Load role info;
	roleId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of roles.html?act=edit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	sql := "SELECT id, role_name, role_codes FROM core_roles WHERE id=?"
	row, err := core.DB.Query(sql, roleId)
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query SQL", err)
		}
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&Html.Role.RoleId, &Html.Role.RoleName, &Html.Role.RoleCodes)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusNotFound)
			return
		}
	}
	r.ParseForm()
	Html.Role.RoleName = r.Form.Get("rolename")
	Html.Role.RoleCodes = strings.Join(r.Form["permissions"], ",")
	err = this.updateRole(&Html.Role)
	if err != nil {
		Html.Info.Show = true
		Html.Info.Message = "Faile to save the role"
		Html.Title = "Edit Role"
		Html.Act = "edit"
		err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "role_manager_editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	} else {
		http.Redirect(w, r, "/roles.html", http.StatusFound)
	}
}

func (this RolesApp) del(w http.ResponseWriter, r *http.Request) {
	roleId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of roles.html?act=edit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	sql := "DELETE FROM core_roles WHERE id = ? "
	_, err = core.DB.Exec(sql, roleId)
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query SQL", err)
		}
	}
	http.Redirect(w, r, "/roles.html", http.StatusFound)
}

func (this RolesApp) updateRole(role *roleRow) error {
	if role == nil {
		return errors.New("Illegal call")
	}
	sql := ""
	var err error = nil
	if role.RoleId == 0 {
		sql = "INSERT INTO core_roles (role_name,role_codes) VALUES (?,?)"
		_, err = core.DB.Query(sql, role.RoleName, role.RoleCodes)
	} else {
		sql = "UPDATE core_roles SET role_name=?, role_codes = ? WHERE id = ? "
		_, err = core.DB.Query(sql, role.RoleName, role.RoleCodes, role.RoleId)
	}
	return err
}
