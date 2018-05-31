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
	"strings"
	"tinybi/core"
	"tinybi/model"
	"tinybi/web"

	"github.com/chai2010/gettext-go/gettext"
)

func RoleListPage(w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	if web.AclRedirect(w, r, "ROLE_ADMIN_VIEW", "/login.html") {
		err := web.GetTemplate(lang, "core/role_manager.html").Execute(w, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

func RoleList(w http.ResponseWriter, r *http.Request) {
	if web.AclRedirect(w, r, "ROLE_ADMIN_VIEW", "/login.html") {
		w.Header().Set("Content-Type", "application/json")
		nullRet := `{"data":[]}`
		var fullRet struct {
			Data []struct {
				Id              int64  `json:"0"`
				RoleName        string `json:"1"`
				RoleCodes       string `json:"2"`
				NumberOfMembers int    `json:"3"`
				EditUrl         string `json:"4"`
				DeleteUrl       string `json:"5"`
			} `json:"data"`
		}
		err := core.DB.Table("core_roles").Select(
			"core_roles.id,core_roles.role_name,core_roles.role_codes,ru.number_of_members,'Edit' as edit_url,'Delete' as delete_url").Join(
			"LEFT", "(select role_id,count(id) as number_of_members from core_users group by role_id ) ru",
			"core_roles.id=ru.role_id").Where(
			"1=1").OrderBy("core_roles.id").Find(&fullRet.Data)
		if err != nil {
			log.Println("Fail to get Role list", err)
			w.Write([]byte(nullRet))
			return
		}
		for i, _ := range fullRet.Data {
			editI18n := gettext.Gettext("Edit")
			deleteI18n := gettext.Gettext("Delete")
			fullRet.Data[i].EditUrl = fmt.Sprintf(`<a href="/acl.html?act=roleEdit&id=%d">%s</a>`, fullRet.Data[i].Id, editI18n)
			fullRet.Data[i].DeleteUrl = fmt.Sprintf(`<a href="/acl.html?act=roleDelete&id=%d">%s</a>`, fullRet.Data[i].Id, deleteI18n)
		}
		w.Write([]byte(web.JsonEncode(fullRet)))
	}
}

func RoleEditPage(act string, w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	if web.AclRedirect(w, r, "ROLE_ADMIN_EDIT", "/login.html") {
		var Html struct {
			Title string
			Act   string
			Role  model.CoreRole
			Info  struct {
				Show    bool
				Type    string
				Message string
			}
			Acls []struct {
				Title string
				Acls  []model.CoreAclCode
			}
			CheckedAcls []string
		}
		var aclTitles []string
		err := core.DB.Table("core_acl_codes").Distinct("title").Find(&aclTitles)
		if err == nil {
			for _, title := range aclTitles {
				var tAcls []model.CoreAclCode
				err := core.DB.Table("core_acl_codes").Where("title=?", title).Find(&tAcls)
				if err != nil {
					log.Println("Fail to get role code from title", title, "Error", err)
					continue
				}
				var tGroup struct {
					Title string
					Acls  []model.CoreAclCode
				}
				tGroup.Title = title
				tGroup.Acls = tAcls
				Html.Acls = append(Html.Acls, tGroup)
			}
		}
		switch act {
		case "add":
			Html.Title = "Create New Role"
			Html.Act = "roleAddExec"
			break
		case "edit":
			Html.Title = "Edit Role"
			Html.Act = "roleEditExec"
			//Load Role info;
			RoleId := r.URL.Query().Get("id")
			if RoleId == "" {
				log.Println("Visit role edit page with null ID")
				web.ErrorNotFound(w, r)
				return
			}
			ok, err := core.DB.Table("core_roles").Where("id=?", RoleId).Get(&Html.Role)
			if !ok {
				log.Println("Visit role edit page with illegal ID")
				if err != nil {
					log.Println(err)
				}
				web.ErrorNotFound(w, r)
				return
			}
			Html.CheckedAcls = strings.Split(Html.Role.RoleCodes, ",")
			break
		}

		err = web.GetTemplate(lang, "core/Role_manager_editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	}
}

func RoleEditExec(act string, w http.ResponseWriter, r *http.Request) {
	lang := web.GetUILang(w, r)
	if web.AclRedirect(w, r, "ROLE_ADMIN_EDIT", "/login.html") {
		r.ParseForm()
		var Html struct {
			Title string
			Act   string
			Role  model.CoreRole
			Info  struct {
				Show    bool
				Type    string
				Message string
			}
			Acls []model.CoreAclCode
		}
		Html.Role.RoleName = r.Form.Get("rolename")
		Html.Role.RoleCodes = strings.Join(r.Form["permissions"], ",")
		switch act {
		case "add":
			Html.Title = "Create New Role"
			Html.Act = "roleAddExec"
			_, err := core.DB.Table("core_roles").Insert(&Html.Role)
			if err != nil {
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = err.Error()
			}
			break
		case "edit":
			Html.Title = "Edit Role"
			Html.Act = "roleEditExec"
			//Load Role info;
			RoleId := r.URL.Query().Get("id")
			if RoleId == "" {
				log.Println("Visit role edit page with null ID")
				web.ErrorNotFound(w, r)
				return
			}
			var role model.CoreRole
			ok, err := core.DB.Table("core_roles").Where("id=?", RoleId).Get(&role)
			if !ok {
				log.Println("Visit role edit page with illegal ID")
				if err != nil {
					log.Println(err)
				}
				web.ErrorNotFound(w, r)
				return
			}
			Html.Role.Id = role.Id
			_, err = core.DB.Table("core_roles").Where("id=?", RoleId).Update(&Html.Role)
			if err != nil {
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = err.Error()
			}
			break
		}

		if Html.Info.Show {
			err := core.DB.Table("core_acl_codes").Where("1=1").Find(&Html.Acls)
			if err != nil {
				log.Println("Fail to get role code list", err)
				web.ErrorNotFound(w, r)
				return
			}

			err = web.GetTemplate(lang, "core/role_manager_editor.html").Execute(w, Html)
			if err != nil {
				log.Println(err)
			}
		} else {
			http.Redirect(w, r, "/acl.html?act=roles", http.StatusFound)
		}
	}
}

func RoleDeleteExec(w http.ResponseWriter, r *http.Request) {
	if web.AclRedirect(w, r, "ROLE_ADMIN_EDIT", "/login.html") {
		var Role model.CoreRole
		RoleId := r.URL.Query().Get("id")
		var user model.CoreUser
		NumofMembers, _ := core.DB.Table("core_users").Where("role_id=?", RoleId).Count(&user)
		if NumofMembers > 0 {
			http.Redirect(w, r, "/acl.html?act=roles", http.StatusFound)
			return
		}
		_, err := core.DB.Table("core_roles").Where("id=?", RoleId).Delete(&Role)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/acl.html?act=roles", http.StatusFound)
	}
}
