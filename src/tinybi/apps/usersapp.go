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
	"log"
	"net/http"
	"strconv"
	"tinybi/core"
	"tinybi/webcore"
)

type UsersApp struct {
	webcore.BaseWebApp
}

type userRow struct {
	Uid      int64  `json:"0"`
	UserName string `json:"1"`
	NickName string `json:"2"`
	RoleName string `json:"3"`
	Status   string `json:"4"`
	EditLink string `json:"5"`
	DelLink  string `json:"6"`
}

type userEntity struct {
	Id               int64
	UserName         string
	NickName         string
	RoleName         string
	Status           string
	RescureCode      string
	Password         string
	RoleId           int64
	Source           string
	LastPasswordTime string
	LastPassword     string
	PasswordSalt     string
}

func (this UsersApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		webcore.AclCheckRedirect(w, r, "SYSTEM", "/login.html")
		switch r.URL.Query().Get("act") {
		case "list":
			//List concurrent tasks which are owned by current user;
			//Super Administrators can see all tasks;
			this.list(w, r)
			break
		case "edit":
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

func (this UsersApp) showPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Roles []struct {
			RoleId   int64
			RoleName string
		}
	}
	sql := "SELECT id, role_name FROM core_roles"
	row, err := core.DB.Query(sql)
	defer row.Close()
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query SQL", err)
		}
	} else {
		for row.Next() {
			var rl struct {
				RoleId   int64
				RoleName string
			}
			err = row.Scan(&rl.RoleId, &rl.RoleName)
			if err == nil {
				Html.Roles = append(Html.Roles, rl)
			}
		}
	}
	err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "user_manager.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this UsersApp) list(w http.ResponseWriter, r *http.Request) {
	//AJAX Method;
	nullRet := `{"data":[]}`
	var fullRet struct {
		Data []userRow `json:"data"`
	}
	//Return JSON Data;
	sql := "SELECT core_users.id AS uid, core_users.user_name, "
	sql += "core_users.nick_name, core_roles.role_name, core_users.status "
	sql += "FROM core_users, core_roles WHERE core_users.role_id = core_roles.id "

	row, err := core.DB.Query(sql)
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query SQL", err)
		}
		w.Write([]byte(nullRet))
		return
	}
	defer row.Close()
	urs := make([]userRow, 0)
	for row.Next() {
		ur := userRow{}
		err = row.Scan(&ur.Uid, &ur.UserName, &ur.NickName, &ur.RoleName, &ur.Status)
		if err != nil {
			w.Write([]byte(nullRet))
			return
		}
		ur.EditLink = fmt.Sprintf("<p class='fa fa-edit'><a href='/users.html?act=edit&id=%d'>Edit</a></p>", ur.Uid)
		ur.DelLink = fmt.Sprintf("<p class='fa fa-trash-o'><a href='/users.html?act=del&id=%d'>Delete</a></p>", ur.Uid)

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

func (this UsersApp) add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := userEntity{}
	user.UserName = r.Form.Get("email")
	user.NickName = r.Form.Get("nickname")
	user.Password = r.Form.Get("password")
	roleId, err := strconv.Atoi(r.Form.Get("role"))
	if err != nil {
		if core.Conf.Debug {
			log.Println("Illegal visit", err)
		}
		http.Redirect(w, r, "/users.html", http.StatusFound)
		return
	}
	user.RoleId = int64(roleId)
	user.PasswordSalt = core.RandomString(64)
	user.Status = r.Form.Get("active")
	if user.Status == "" {
		user.Status = "INACTIVE"
	}
	//Preserved column;
	user.Source = "A"
	err = this.updateUser(&user)
	if core.Conf.Debug {
		println(err)
	}
	http.Redirect(w, r, "/users.html", http.StatusFound)
}

func (this UsersApp) edit(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Illegal visit of users.html?act=edit", err)
		webcore.ErrorNotFound(w, r)
		return
	}
	var Html struct {
		Title string
		Act   string
		User  userEntity
		Info  struct {
			Show    bool
			Type    string
			Message string
		}
		Roles []struct {
			RoleId   int64
			RoleName string
		}
	}
	sql := `SELECT
		id,
		user_name,
		nick_name,
		status,
		rescure_code,
		password,
		role_id,
		source,
		last_password_time,
		last_password,
		password_salt
	FROM
		core_users `
	sql += "WHERE id = ?"
	err = core.DB.QueryRow(sql, userId).Scan(
		&Html.User.Id,
		&Html.User.UserName,
		&Html.User.NickName,
		&Html.User.Status,
		&Html.User.RescureCode,
		&Html.User.Password,
		&Html.User.RoleId,
		&Html.User.Source,
		&Html.User.LastPasswordTime,
		&Html.User.LastPassword,
		&Html.User.PasswordSalt)
	if err != nil {
		log.Println("Illegal visit of users.html?act=edit", err)
		webcore.ErrorNotFound(w, r)
		return
	}
	r.ParseForm()
	Html.User.UserName = r.Form.Get("email")
	Html.User.NickName = r.Form.Get("nickname")
	newPassword := r.Form.Get("password")
	if newPassword != "" {
		Html.User.Password = newPassword
	}
	roleId, err := strconv.Atoi(r.Form.Get("role"))
	if err != nil {
		if core.Conf.Debug {
			log.Println("Illegal visit", err)
		}
		http.Redirect(w, r, "/users.html", http.StatusFound)
		return
	}
	Html.User.RoleId = int64(roleId)
	Html.User.Status = r.Form.Get("active")
	if Html.User.Status == "" {
		Html.User.Status = "INACTIVE"
	}
	//Preserved column;
	Html.User.Source = "A"
	if Html.User.UserName == "" || Html.User.NickName == "" {
		Html.Info.Show = true
		Html.Info.Type = "danger"
		Html.Info.Message = "User name and nick cannot be null"
		Html.Title = "Edit User"
		Html.Act = "edit"
		err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "user_manager_editor.html").Execute(w, Html)
	}
	err = this.updateUser(&Html.User)
	Html.Title = "Edit Role"
	Html.Act = "edit"
	if err != nil {
		if core.Conf.Debug {
			println(err)
		}
		Html.Info.Show = true
		Html.Info.Type = "danger"
		Html.Info.Message = "Faile to save the user"
		Html.Title = "Edit User"
		Html.Act = "edit"
		err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "user_manager_editor.html").Execute(w, Html)
	} else {
		http.Redirect(w, r, "/users.html", http.StatusFound)
	}
}

func (this UsersApp) del(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of users.html?act=del")
		webcore.ErrorNotFound(w, r)
		return
	}
	sql := "DELETE FROM core_users WHERE id = ? "
	_, err = core.DB.Exec(sql, userId)
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query SQL", err)
		}
	}
	http.Redirect(w, r, "/roles.html", http.StatusFound)
}

func (this UsersApp) updateUser(user *userEntity) error {
	if user == nil {
		return errors.New("Illegal call")
	}
	sql := ""
	var err error = nil
	if user.Id == 0 {
		sql = "INSERT INTO core_users "
		sql += `(user_name,
				nick_name,
				status,
				rescure_code,
				password,
				role_id,
				source,
				last_password_time,
				last_password,
				password_salt) VALUES `
		sql += "(?,?,?,?,md5( concat( md5(?), core_users.password_salt )),?,?,?,?,?)"
		_, err = core.DB.Query(sql, user.UserName, user.NickName,
			user.Status, user.RescureCode, user.Password, user.RoleId,
			user.Source, user.LastPasswordTime, user.LastPassword,
			user.PasswordSalt)
	} else {
		sql = "UPDATE core_users SET "
		sql += `user_name=?,
		nick_name=?,
		status=?,
		rescure_code=?,
		password=md5( concat( md5(?), core_users.password_salt )),
		role_id=?,
		source=?,
		last_password_time=?,
		last_password=?,
		password_salt=? `
		sql += "WHERE id = ?"
		_, err = core.DB.Query(sql, user.UserName, user.NickName,
			user.Status, user.RescureCode, user.Password, user.RoleId,
			user.Source, user.LastPasswordTime, user.LastPassword,
			user.PasswordSalt, user.Id)
	}
	return err
}

func (this UsersApp) editPage(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("Illegal visit of users.html?act=edit", err)
		webcore.ErrorNotFound(w, r)
		return
	}
	var Html struct {
		Title string
		Act   string
		User  userEntity
		Info  struct {
			Show    bool
			Type    string
			Message string
		}
		Roles []struct {
			RoleId   int64
			RoleName string
		}
	}
	sql := "SELECT id, role_name FROM core_roles"
	row, err := core.DB.Query(sql)
	defer row.Close()
	if err != nil {
		if core.Conf.Debug {
			log.Println("Fail to query SQL", err)
		}
	} else {
		for row.Next() {
			var rl struct {
				RoleId   int64
				RoleName string
			}
			err = row.Scan(&rl.RoleId, &rl.RoleName)
			if err == nil {
				Html.Roles = append(Html.Roles, rl)
			}
		}
	}
	sql = `SELECT
		id,
		user_name,
		nick_name,
		status,
		rescure_code,
		password,
		role_id,
		source,
		last_password_time,
		last_password,
		password_salt
	FROM
		core_users `
	sql += "WHERE id = ?"
	err = core.DB.QueryRow(sql, userId).Scan(
		&Html.User.Id,
		&Html.User.UserName,
		&Html.User.NickName,
		&Html.User.Status,
		&Html.User.RescureCode,
		&Html.User.Password,
		&Html.User.RoleId,
		&Html.User.Source,
		&Html.User.LastPasswordTime,
		&Html.User.LastPassword,
		&Html.User.PasswordSalt)
	if err != nil {
		log.Println("Illegal visit of users.html?act=edit", err)
		webcore.ErrorNotFound(w, r)
		return
	}
	Html.Title = "Edit user"
	Html.Act = "edit"
	err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "user_manager_editor.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}
