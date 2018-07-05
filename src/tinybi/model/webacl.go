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
package model

import (
	"log"
	"time"
	"tinybi/core"
	"crypto/md5"
	"fmt"
	"io"
)

//Access Control Control Models for WEB access;

type CoreUser struct {
	Id               int64
	EMail            string    `xorm:"notnull unique 'email'"`
	NickName         string    `xorm:"notnull unique 'nick_name'"`
	Source           string    `xorm:"'source' default 'A'"`
	Lang             string    `xorm:"'lang'"`
	Status           string    `xorm:"'status' default 'DISABLED'"`
	RoleId           int64     `xorm:"'role_id'"`
	RescureCode      string    `xorm:"'rescure_code'"`
	Password         string    `xorm:"'password'"`
	PasswordSalt     string    `xorm:"'password_salt'"`
	LastPasswordTime time.Time `xorm:"last_password_time"`
	LastPassword     string    `xorm:"last_password"`
	LastUpdated      time.Time `xorm:"TIMESTAMP 'last_updated' created updated default CURRENT_TIMESTAMP"`
}

type CoreRole struct {
	Id          int64
	RoleName    string    `xorm:"'role_name'"`
	RoleCodes   string    `xorm:"TEXT 'role_codes'"`
	LastUpdated time.Time `xorm:"TIMESTAMP 'last_updated' created updated default CURRENT_TIMESTAMP"`
}

type CoreAclCode struct {
	Id          int64
	Title       string `xorm:"'title'"`
	Code        string `xorm:"'code'"`
	Description string `xorm:"'Description'"`
}

func CoreUserPassword(rawPassword string, salt string) string {
	if rawPassword == "" || salt == "" {
		return ""
	}
	pHash := md5.New()
	io.WriteString(pHash, rawPassword)
	pp := fmt.Sprintf("%x", pHash.Sum(nil))
	pp += salt
	fpHash := md5.New()
	io.WriteString(fpHash, pp)
	fp := fmt.Sprintf("%x", fpHash.Sum(nil))
	return fp
}

func InstallAclTables() {
	core.DB.Table("core_users").Sync2(new(CoreUser))
	core.DB.Table("core_roles").Sync2(new(CoreRole))
	core.DB.Table("core_acl_codes").Sync2(new(CoreAclCode))
	//Insert basic acl codes if they are not exist;
	basicCodes := make([][]string, 4)
	basicCodes[0] = []string{"User Mangement", "USER_ADMIN_VIEW", "User Management, Readonly"}
	basicCodes[1] = []string{"User Mangement", "USER_ADMIN_EDIT", "Edit & Delete Users"}
	basicCodes[2] = []string{"Role Mangement", "ROLE_ADMIN_VIEW", "Role Management, Readonly"}
	basicCodes[3] = []string{"Role Mangement", "ROLE_ADMIN_EDIT", "Edit & Delete Roles"}
	for _, basicCode := range basicCodes {
		var code CoreAclCode
		ok, err := core.DB.Table("core_acl_codes").Where("code=?", basicCode[1]).Get(&code)
		if ok {
			continue
		}
		if err != nil {
			log.Println("Error occurs when initializing ACL codes", err)
		}
		code.Title = basicCode[0]
		code.Code = basicCode[1]
		code.Description = basicCode[2]
		_, err = core.DB.Table("core_acl_codes").Insert(&code)
		if err != nil {
			log.Println("Error occurs when inserting ACL code", code.Code, "Error", err)
		}
	}
}
