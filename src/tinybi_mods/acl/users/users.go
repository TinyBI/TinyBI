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
package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tinybi/core"
	"tinybi/model"
	"tinybi/web"

	"github.com/chai2010/gettext-go/gettext"
)

func UserListPage(w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	if web.AclRedirect(w, r, "USER_ADMIN_VIEW", "/login.html") {
		err := web.GetTemplate(lang, "core/user_manager.html").Execute(w, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

func UserList(w http.ResponseWriter, r *http.Request) {
	if web.AclRedirect(w, r, "USER_ADMIN_VIEW", "/login.html") {
		w.Header().Set("Content-Type", "application/json")
		nullRet := `{"data":[]}`
		var fullRet struct {
			Data []struct {
				Id        int64  `json:"0"`
				Email     string `json:"1"`
				NickName  string `json:"2"`
				RoleName  string `json:"3"`
				Status    string `json:"4"`
				EditUrl   string `json:"5"`
				DeleteUrl string `json:"6"`
			} `json:"data"`
		}
		err := core.DB.Table("core_users").Select(
			"core_users.id,core_users.email,core_users.nick_name,core_roles.role_name,core_users.status,'Edit' as edit_url,'Delete' as delete_url").Join(
			"INNER", "core_roles", "core_users.role_id=core_roles.id").Where(
			"1=1").OrderBy("core_users.id").Find(&fullRet.Data)
		if err != nil {
			log.Println("Fail to get user list", err)
			w.Write([]byte(nullRet))
			return
		}
		for i := range fullRet.Data {
			editI18n := gettext.Gettext("Edit")
			deleteI18n := gettext.Gettext("Delete")
			fullRet.Data[i].EditUrl = fmt.Sprintf(`<a href="/acl.html?act=userEdit&id=%d">%s</a>`, fullRet.Data[i].Id, editI18n)
			fullRet.Data[i].DeleteUrl = fmt.Sprintf(`<a href="/acl.html?act=userDelete&id=%d">%s</a>`, fullRet.Data[i].Id, deleteI18n)
		}
		w.Write([]byte(web.JsonEncode(fullRet)))
	}
}

func UserEditPage(act string, w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	if web.AclRedirect(w, r, "USER_ADMIN_EDIT", "/login.html") {
		var Html struct {
			Title string
			Act   string
			User  model.CoreUser
			Info  struct {
				Show    bool
				Type    string
				Message string
			}
			Roles []model.CoreRole
		}
		err := core.DB.Table("core_roles").Where("1=1").Find(&Html.Roles)
		if err != nil {
			log.Println("Fail to get role list", err)
			web.ErrorNotFound(w, r)
			return
		}
		switch act {
		case "add":
			Html.Title = "Create New User"
			Html.Act = "userAddExec"
			break
		case "edit":
			Html.Title = "Edit User"
			Html.Act = "userEditExec"
			//Load User info;
			userId := r.URL.Query().Get("id")
			if userId == "" {
				log.Println("Visit user edit page with null ID")
				web.ErrorNotFound(w, r)
				return
			}
			ok, err := core.DB.Table("core_users").Where("id=?", userId).Get(&Html.User)
			if !ok {
				log.Println("Visit user edit page with illegal ID")
				if err != nil {
					log.Println(err)
				}
				web.ErrorNotFound(w, r)
				return
			}
			break
		}

		err = web.GetTemplate(lang, "core/user_manager_editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	}
}

func UserEditExec(act string, w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	if web.AclRedirect(w, r, "USER_ADMIN_EDIT", "/login.html") {
		r.ParseForm()
		var Html struct {
			Title string
			Act   string
			User  model.CoreUser
			Info  struct {
				Show    bool
				Type    string
				Message string
			}
			Roles []model.CoreRole
		}
		Html.User.EMail = r.Form.Get("email")
		Html.User.NickName = r.Form.Get("nickname")
		Html.User.Password = r.Form.Get("password")
		roleId, err := strconv.Atoi(r.Form.Get("role"))
		if err != nil {
			if core.Conf.Debug {
				log.Println("Illegal visit", err)
			}
			web.ErrorNotFound(w, r)
			return
		}
		Html.User.RoleId = int64(roleId)
		Html.User.PasswordSalt = core.RandomString(64)
		Html.User.Password = model.CoreUserPassword(Html.User.Password, Html.User.PasswordSalt)
		Html.User.Status = r.Form.Get("active")
		if Html.User.Status == "" {
			Html.User.Status = "INACTIVE"
		}
		//Preserved column;
		Html.User.Source = "A"
		switch act {
		case "add":
			Html.Title = "Create New User"
			Html.Act = "userAddExec"
			_, err = core.DB.Table("core_users").Insert(&Html.User)
			if err != nil {
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = err.Error()
			}
			break
		case "edit":
			Html.Title = "Edit User"
			Html.Act = "userEditExec"
			//Load User info;
			userId := r.URL.Query().Get("id")
			if userId == "" {
				log.Println("Visit user edit page with null ID")
				web.ErrorNotFound(w, r)
				return
			}
			var user model.CoreUser
			ok, err := core.DB.Table("core_users").Where("id=?", userId).Get(&user)
			if !ok {
				log.Println("Visit user edit page with illegal ID")
				if err != nil {
					log.Println(err)
				}
				web.ErrorNotFound(w, r)
				return
			}
			Html.User.Id = user.Id
			Html.User.PasswordSalt = user.PasswordSalt
			if Html.User.Password != "" {
				//Update password
				newPassword := fmt.Sprintf("md5( concat( md5(%s), password_salt ))", Html.User.Password)
				Html.User.Password = newPassword
			}
			_, err = core.DB.Table("core_users").Where("id=?", userId).Update(&Html.User)
			if err != nil {
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = err.Error()
			}
			break
		}

		if Html.Info.Show {
			err = core.DB.Table("core_roles").Where("1=1").Find(&Html.Roles)
			if err != nil {
				log.Println("Fail to get role list", err)
				web.ErrorNotFound(w, r)
				return
			}

			err = web.GetTemplate(lang, "core/user_manager_editor.html").Execute(w, Html)
			if err != nil {
				log.Println(err)
			}
		} else {
			http.Redirect(w, r, "/acl.html?act=users", http.StatusFound)
		}
	}
}

func UserDeleteExec(w http.ResponseWriter, r *http.Request) {
	if web.AclRedirect(w, r, "USER_ADMIN_EDIT", "/login.html") {
		var user model.CoreUser
		userId := r.URL.Query().Get("id")
		_, err := core.DB.Table("core_users").Where("id=?", userId).Delete(&user)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/acl.html?act=users", http.StatusFound)
	}
}
