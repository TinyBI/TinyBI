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
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"tinybi/core"
	"tinybi/models"
	"tinybi/webcore"
)

type GLApp struct {
	webcore.BaseWebApp
}

func (this GLApp) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		switch r.URL.Query().Get("act") {
		case "periods":
			//Show page of accounting periods;
			webcore.AclCheckRedirect(w, r, "GL_PERIODS_R", "/login.html")
			this.periodPage(w, r)
			break
		case "periodsList":
			//AJAX Method, list all accounting periods;
			webcore.AclCheckRedirect(w, r, "GL_PERIODS_R", "/login.html")
			this.periodList(w, r)
			break
		case "periodAdd":
			webcore.AclCheckRedirect(w, r, "GL_PERIODS_W", "/login.html")
			this.periodAddPage(w, r)
			break
		case "periodEdit":
			webcore.AclCheckRedirect(w, r, "GL_PERIODS_W", "/login.html")
			this.periodEditPage(w, r)
			break
		case "periodOpen":
			webcore.AclCheckRedirect(w, r, "GL_PERIODS_W", "/login.html")
			this.periodOpen(w, r)
			break
		case "accounts":
			webcore.AclCheckRedirect(w, r, "GL_ACCOUNTS_R", "/login.html")
			this.accountPage(w, r)
			break
		case "accountsList":
			//AJAX Method, list all accounts;
			webcore.AclCheckRedirect(w, r, "GL_ACCOUNTS_R", "/login.html")
			this.accountList(w, r)
			break
		case "accountAdd":
			webcore.AclCheckRedirect(w, r, "GL_ACCOUNTS_W", "/login.html")
			this.accountAddPage(w, r)
			break
		case "accountEdit":
			webcore.AclCheckRedirect(w, r, "GL_ACCOUNTS_W", "/login.html")
			this.accountEditPage(w, r)
			break
		case "sobs":
			webcore.AclCheckRedirect(w, r, "GL_SOBS_R", "/login.html")
			this.sobPage(w, r)
			break
		case "sobsList":
			//AJAX Method, list all accounts;
			webcore.AclCheckRedirect(w, r, "GL_SOBS_R", "/login.html")
			this.sobList(w, r)
			break
		case "sobAdd":
			webcore.AclCheckRedirect(w, r, "GL_SOBS_W", "/login.html")
			this.sobAddPage(w, r)
			break
		case "sobEdit":
			webcore.AclCheckRedirect(w, r, "GL_SOBS_W", "/login.html")
			this.sobEditPage(w, r)
			break
		case "journals":
			webcore.AclCheckRedirect(w, r, "GL_JES_R", "/login.html")
			this.journalPage(w, r)
			break
		case "journalAdd":
			webcore.AclCheckRedirect(w, r, "GL_JES_CREATE", "/login.html")
			this.journalAddPage(w, r)
			break
		default:
			//404
			http.Redirect(w, r, "/", http.StatusNotFound)
			return
		}
	} else {
		switch r.URL.Query().Get("act") {
		case "periodAdd":
			webcore.AclCheckRedirect(w, r, "GL_PERIODS_W", "/login.html")
			this.periodAdd(w, r)
			break
		case "periodEdit":
			webcore.AclCheckRedirect(w, r, "GL_PERIODS_W", "/login.html")
			this.periodEdit(w, r)
			break
		case "accountAdd":
			webcore.AclCheckRedirect(w, r, "GL_ACCOUNTS_W", "/login.html")
			this.accountAdd(w, r)
			break
		case "accountEdit":
			webcore.AclCheckRedirect(w, r, "GL_ACCOUNTS_W", "/login.html")
			this.accountEdit(w, r)
			break
		case "sobAdd":
			webcore.AclCheckRedirect(w, r, "GL_SOBS_W", "/login.html")
			this.sobAdd(w, r)
			break
		case "sobEdit":
			webcore.AclCheckRedirect(w, r, "GL_SOBS_W", "/login.html")
			this.sobEdit(w, r)
			break
		case "journalAdd":
			webcore.AclCheckRedirect(w, r, "GL_JES_CREATE", "/login.html")
			this.journalAdd(w, r)
			break
		}
	}
}

func (this GLApp) periodPage(w http.ResponseWriter, r *http.Request) {
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_periods.html").Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) periodList(w http.ResponseWriter, r *http.Request) {
	nullRet := `{"data":[]}`
	var fullRet struct {
		Data []struct {
			PeriodCode  string `json:"0"`
			PeriodName  string `json:"1"`
			Status      string `json:"2"`
			Description string `json:"3"`
			StartTime   string `json:"4"`
			EndTime     string `json:"5"`
			EditLink    string `json:"6"`
		} `json:"data"`
	}
	var periods []models.GLPeriod
	err := core.DBEngine.Table("gl_periods").Find(&periods)
	if err != nil {
		w.Write([]byte(nullRet))
		return
	}
	for _, period := range periods {
		var periodRow struct {
			PeriodCode  string `json:"0"`
			PeriodName  string `json:"1"`
			Status      string `json:"2"`
			Description string `json:"3"`
			StartTime   string `json:"4"`
			EndTime     string `json:"5"`
			EditLink    string `json:"6"`
		}
		periodRow.PeriodCode = period.PeriodCode
		periodRow.PeriodName = period.PeriodName
		periodRow.Status = period.Status
		periodRow.Description = period.Description
		periodRow.StartTime = time.Unix(int64(period.StartTime), 0).Format("2006-01-02 15:04:05")
		periodRow.EndTime = time.Unix(int64(period.EndTime), 0).Format("2006-01-02 15:04:05")
		//periodRow.EditLink = fmt.Sprintf("<p class='fa fa-edit'><a href='/gl.html?act=periodEdit&id=%d'>Edit</a></P>", period.Id)
		checked := ""
		if period.Status == "OPENED" {
			checked = "checked"
		}
		periodRow.EditLink = fmt.Sprintf("<input class='periodStatus' type='checkbox' data-toggle='toggle' id='period%d' onchange=\"openPeriod('#period%d',%d)\" %s >", period.Id, period.Id, period.Id, checked)
		fullRet.Data = append(fullRet.Data, periodRow)
	}
	sret, err := json.Marshal(fullRet)
	if err != nil {
		w.Write([]byte(nullRet))
		return
	}
	w.Write(sret)
}

func (this GLApp) periodAddPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title  string
		Period models.GLPeriod
		Act    string
		Info   struct {
			Show    bool
			Type    string
			Message string
		}
	}
	Html.Title = "Add Period"
	Html.Act = "periodAdd"
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_periods_editor.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) periodAdd(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title  string
		Period models.GLPeriod
		Act    string
		Info   struct {
			Show    bool
			Type    string
			Message string
		}
	}
	r.ParseForm()
	Html.Period.PeriodCode = r.Form.Get("periodcode")
	Html.Period.PeriodName = r.Form.Get("periodname")
	Html.Period.Status = r.Form.Get("status")
	Html.Period.Description = r.Form.Get("description")
	Html.Period.StartTime = core.UnixTime(r.Form.Get("starttime"))
	Html.Period.EndTime = core.UnixTime(r.Form.Get("endtime"))
	Html.Title = "Add Period"
	Html.Act = "periodAdd"
	if Html.Period.PeriodCode == "" || Html.Period.PeriodName == "" {
		Html.Info.Show = true
		Html.Info.Type = "danger"
		Html.Info.Message = "You must input the period code and name"
	} else {
		if Html.Period.StartTime == 0 || Html.Period.EndTime == 0 {
			Html.Info.Show = true
			Html.Info.Type = "danger"
			Html.Info.Message = "You must select the time span of this period"
		}
	}
	Html.Period.LastUpdated = time.Now()
	if !Html.Info.Show {
		//Check conflict period;
		var dupPeriod models.GLPeriod
		ok, err := core.DBEngine.Table("gl_periods").Where("start_time<=?", Html.Period.EndTime).And("end_time>=?", Html.Period.StartTime).And("id!=?", Html.Period.Id).Get(&dupPeriod)
		if ok {
			Html.Info.Show = true
			Html.Info.Type = "danger"
			Html.Info.Message = "Found conflict period: "
			Html.Info.Message += dupPeriod.PeriodName
			Html.Info.Message += ", time span: "
			Html.Info.Message += time.Unix(int64(dupPeriod.StartTime), 0).Format("2006-01-02 15:04:05")
			Html.Info.Message += " to "
			Html.Info.Message += time.Unix(int64(dupPeriod.EndTime), 0).Format("2006-01-02 15:04:05")
		} else {
			_, err = core.DBEngine.Table("gl_periods").Insert(&Html.Period)
			if err != nil {
				if core.Conf.Debug {
					log.Println(err)
				}
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = "Fail to save the period"
			}
		}
	}
	if Html.Info.Show {
		err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_periods_editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	} else {
		this.periodPage(w, r)
	}
}

func (this GLApp) periodEditPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title  string
		Period models.GLPeriod
		Act    string
		Info   struct {
			Show    bool
			Type    string
			Message string
		}
	}
	//Load period info;
	periodId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of gl.html?act=periodEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	ok, err := core.DBEngine.Table("gl_periods").Where("id = ?", periodId).Get(&Html.Period)
	if !ok {
		if err != nil {
			log.Println(err)
		}
		log.Printf("Illegal visit of gl.html?act=periodEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	Html.Title = "Edit Period"
	Html.Act = "periodEdit"
	err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_periods_editor.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) periodEdit(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title  string
		Period models.GLPeriod
		Act    string
		Info   struct {
			Show    bool
			Type    string
			Message string
		}
	}
	//Load period info;
	periodId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of gl.html?act=periodEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	ok, err := core.DBEngine.Table("gl_periods").Where("id = ?", periodId).Get(&Html.Period)
	if !ok {
		if err != nil {
			log.Println(err)
		}
		log.Printf("Illegal visit of gl.html?act=periodEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	r.ParseForm()
	Html.Period.PeriodCode = r.Form.Get("periodcode")
	Html.Period.PeriodName = r.Form.Get("periodname")
	Html.Period.Status = r.Form.Get("status")
	Html.Period.Description = r.Form.Get("description")
	Html.Period.StartTime = core.UnixTime(r.Form.Get("starttime"))
	Html.Period.EndTime = core.UnixTime(r.Form.Get("endtime"))
	Html.Title = "Edit Period"
	Html.Act = "periodEdit"
	if Html.Period.PeriodCode == "" || Html.Period.PeriodName == "" {
		Html.Info.Show = true
		Html.Info.Type = "danger"
		Html.Info.Message = "You must input the period code and name"
	} else {
		if Html.Period.StartTime == 0 || Html.Period.EndTime == 0 {
			Html.Info.Show = true
			Html.Info.Type = "danger"
			Html.Info.Message = "You must select the time span of this period"
		}
	}
	Html.Period.LastUpdated = time.Now()
	if !Html.Info.Show {
		//Check conflict period;
		var dupPeriod models.GLPeriod
		ok, err := core.DBEngine.Table("gl_periods").Where("start_time<=?", Html.Period.EndTime).And("end_time>=?", Html.Period.StartTime).And("id!=?", Html.Period.Id).Get(&dupPeriod)
		if ok {
			Html.Info.Show = true
			Html.Info.Type = "danger"
			Html.Info.Message = "Found conflict period: "
			Html.Info.Message += dupPeriod.PeriodName
			Html.Info.Message += ", time span: "
			Html.Info.Message += time.Unix(int64(dupPeriod.StartTime), 0).Format("2006-01-02 15:04:05")
			Html.Info.Message += " to "
			Html.Info.Message += time.Unix(int64(dupPeriod.EndTime), 0).Format("2006-01-02 15:04:05")
		} else {
			if err != nil && core.Conf.Debug {
				log.Println(err)
			}
			_, err = core.DBEngine.Table("gl_periods").Where("id = ?", periodId).Update(&Html.Period)
			if err != nil {
				if core.Conf.Debug {
					log.Println(err)
				}
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = "Fail to save the period"
			}
		}
	}
	if Html.Info.Show {
		err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_periods_editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	} else {
		this.periodPage(w, r)
	}
}

func (this GLApp) periodOpen(w http.ResponseWriter, r *http.Request) {
	periodId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of gl.html?act=periodOpen")
		w.Write([]byte("0"))
		return
	}
	var period models.GLPeriod
	ok, err := core.DBEngine.Table("gl_periods").Where("id = ?", periodId).Get(&period)
	if !ok {
		if err != nil {
			log.Println(err)
		}
		log.Printf("Illegal visit of gl.html?act=periodOpen")
		w.Write([]byte("0"))
		return
	}
	if period.Status == "OPENED" {
		period.Status = "CLOSED"
	} else {
		period.Status = "OPENED"
	}
	_, err = core.DBEngine.Table("gl_periods").Where("id = ?", periodId).Update(&period)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
		}
		w.Write([]byte("0"))
		return
	}
	w.Write([]byte("1"))
}

func (this GLApp) accountPage(w http.ResponseWriter, r *http.Request) {
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_accounts.html").Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) accountList(w http.ResponseWriter, r *http.Request) {
	nullRet := `{"data":[]}`
	var fullRet struct {
		Data []struct {
			AccountCode string `json:"0"`
			AccountName string `json:"1"`
			Description string `json:"2"`
			EditLink    string `json:"3"`
		} `json:"data"`
	}
	var accounts []models.GLAccount
	err := core.DBEngine.Table("gl_accounts").Find(&accounts)
	if err != nil {
		w.Write([]byte(nullRet))
		return
	}
	for _, account := range accounts {
		var accountRow struct {
			AccountCode string `json:"0"`
			AccountName string `json:"1"`
			Description string `json:"2"`
			EditLink    string `json:"3"`
		}
		accountRow.AccountCode = account.AccountCode
		accountRow.AccountName = account.AccountName
		accountRow.Description = account.Description
		accountRow.EditLink = fmt.Sprintf("<p class='fa fa-edit'><a href='/gl.html?act=accountEdit&id=%d'>Edit</a></P>", account.Id)
		fullRet.Data = append(fullRet.Data, accountRow)
	}
	sret, err := json.Marshal(fullRet)
	if err != nil {
		w.Write([]byte(nullRet))
		return
	}
	w.Write(sret)
}

func (this GLApp) accountAddPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title   string
		Account models.GLAccount
		Act     string
		Info    struct {
			Show    bool
			Type    string
			Message string
		}
	}
	Html.Title = "Add Account"
	Html.Act = "accountAdd"
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_accounts_editor.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) accountAdd(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title   string
		Account models.GLAccount
		Act     string
		Info    struct {
			Show    bool
			Type    string
			Message string
		}
	}
	r.ParseForm()
	Html.Account.AccountCode = r.Form.Get("accountcode")
	Html.Account.AccountName = r.Form.Get("accountname")
	Html.Account.Description = r.Form.Get("description")
	Html.Title = "Add Account"
	Html.Act = "accountAdd"
	if Html.Account.AccountCode == "" || Html.Account.AccountName == "" {
		Html.Info.Show = true
		Html.Info.Type = "danger"
		Html.Info.Message = "You must input the account code and name"
	}
	Html.Account.LastUpdated = time.Now()
	if !Html.Info.Show {
		//Check conflict account;
		var dupAccount models.GLAccount
		ok, err := core.DBEngine.Table("gl_accounts").Where("account_code=?", Html.Account.AccountCode).And("id!=?", Html.Account.Id).Get(&dupAccount)
		if ok {
			Html.Info.Show = true
			Html.Info.Type = "danger"
			Html.Info.Message = "Found conflict account: "
			Html.Info.Message += dupAccount.AccountName
			Html.Info.Message += " ,code: "
			Html.Info.Message += dupAccount.AccountCode
		} else {
			_, err = core.DBEngine.Table("gl_accounts").Insert(&Html.Account)
			if err != nil {
				if core.Conf.Debug {
					log.Println(err)
				}
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = "Fail to save the account"
			}
		}
	}
	if Html.Info.Show {
		err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_accounts_editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	} else {
		this.accountPage(w, r)
	}
}

func (this GLApp) accountEditPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title   string
		Account models.GLAccount
		Act     string
		Info    struct {
			Show    bool
			Type    string
			Message string
		}
	}
	//Load account info;
	accountId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of gl.html?act=accountEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	ok, err := core.DBEngine.Table("gl_accounts").Where("id = ?", accountId).Get(&Html.Account)
	if !ok {
		if err != nil {
			log.Println(err)
		}
		log.Printf("Illegal visit of gl.html?act=accountEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	Html.Title = "Edit Account"
	Html.Act = "accountEdit"
	err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_accounts_editor.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) accountEdit(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title   string
		Account models.GLAccount
		Act     string
		Info    struct {
			Show    bool
			Type    string
			Message string
		}
	}
	//Load account info;
	accountId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of gl.html?act=accountEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	ok, err := core.DBEngine.Table("gl_accounts").Where("id = ?", accountId).Get(&Html.Account)
	if !ok {
		if err != nil {
			log.Println(err)
		}
		log.Printf("Illegal visit of gl.html?act=accountEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	r.ParseForm()
	Html.Account.AccountCode = r.Form.Get("accountcode")
	Html.Account.AccountName = r.Form.Get("accountname")
	Html.Account.Description = r.Form.Get("description")
	Html.Title = "Edit Accounts"
	Html.Act = "accountEdit"
	if Html.Account.AccountCode == "" || Html.Account.AccountName == "" {
		Html.Info.Show = true
		Html.Info.Type = "danger"
		Html.Info.Message = "You must input the account code and name"
	}
	Html.Account.LastUpdated = time.Now()
	if !Html.Info.Show {
		//Check conflict account;
		var dupAccount models.GLAccount
		ok, err := core.DBEngine.Table("gl_accounts").Where("account_code=?", Html.Account.AccountCode).And("id!=?", Html.Account.Id).Get(&dupAccount)
		if ok {
			Html.Info.Show = true
			Html.Info.Type = "danger"
			Html.Info.Message = "Found conflict account: "
			Html.Info.Message += dupAccount.AccountName
			Html.Info.Message += " ,code: "
			Html.Info.Message += dupAccount.AccountCode
		} else {
			if err != nil && core.Conf.Debug {
				log.Println(err)
			}
			_, err = core.DBEngine.Table("gl_accounts").Where("id = ?", accountId).Update(&Html.Account)
			if err != nil {
				if core.Conf.Debug {
					log.Println(err)
				}
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = "Fail to save the account"
			}
		}
	}
	if Html.Info.Show {
		err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_accounts_editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	} else {
		this.accountPage(w, r)
	}
}

func (this GLApp) sobPage(w http.ResponseWriter, r *http.Request) {
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_sobs.html").Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) sobList(w http.ResponseWriter, r *http.Request) {
	nullRet := `{"data":[]}`
	var fullRet struct {
		Data []struct {
			SobName      string `json:"0"`
			CurrencyCode string `json:"1"`
			EolTime      string `json:"2"`
			EditLink     string `json:"3"`
		} `json:"data"`
	}
	var sobs []models.GLSetOfBook
	err := core.DBEngine.Table("gl_set_of_books").Find(&sobs)
	if err != nil {
		w.Write([]byte(nullRet))
		return
	}
	if len(sobs) == 0 {
		w.Write([]byte(nullRet))
		return
	}
	for _, sob := range sobs {
		var sobRow struct {
			SobName      string `json:"0"`
			CurrencyCode string `json:"1"`
			EolTime      string `json:"2"`
			EditLink     string `json:"3"`
		}
		sobRow.SobName = sob.SobName
		sobRow.CurrencyCode = sob.CurrencyCode
		sobRow.EolTime = core.FromUnixTime(sob.EolTime)
		sobRow.EditLink = fmt.Sprintf("<p class='fa fa-edit'><a href='/gl.html?act=sobEdit&id=%d'>Edit</a></P>", sob.Id)
		fullRet.Data = append(fullRet.Data, sobRow)
	}
	sret, err := json.Marshal(fullRet)
	if err != nil {
		w.Write([]byte(nullRet))
		return
	}
	w.Write(sret)
}

func (this GLApp) sobAddPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title    string
		Sob      models.GLSetOfBook
		Act      string
		Timezone string
		Info     struct {
			Show    bool
			Type    string
			Message string
		}
	}
	Html.Title = "Create set of book"
	Html.Act = "sobAdd"
	Html.Timezone = core.Conf.TimeZone
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_sobs_editor.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) sobAdd(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title    string
		Sob      models.GLSetOfBook
		Timezone string
		Act      string
		Info     struct {
			Show    bool
			Type    string
			Message string
		}
	}
	Html.Timezone = core.Conf.TimeZone
	r.ParseForm()
	Html.Sob.SobName = r.Form.Get("sobname")
	Html.Sob.CurrencyCode = r.Form.Get("currencycode")
	Html.Sob.EolTime = core.UnixTime(r.Form.Get("eoltime"))
	Html.Title = "Create set of book"
	Html.Act = "sobAdd"
	if Html.Sob.SobName == "" || Html.Sob.CurrencyCode == "" || Html.Sob.EolTime == 0 {
		Html.Info.Show = true
		Html.Info.Type = "danger"
		Html.Info.Message = "All input fields are required"
	}
	Html.Sob.LastUpdated = time.Now()
	if !Html.Info.Show {
		//Check conflict account;
		var dupSob models.GLSetOfBook
		ok, err := core.DBEngine.Table("gl_set_of_books").Where("sob_name=?", Html.Sob.SobName).And("id!=?", Html.Sob.Id).Get(&dupSob)
		if ok {
			Html.Info.Show = true
			Html.Info.Type = "danger"
			Html.Info.Message = "Found conflict account: "
			Html.Info.Message += dupSob.SobName
			Html.Info.Message += " ,currency code: "
			Html.Info.Message += dupSob.CurrencyCode
		} else {
			_, err = core.DBEngine.Table("gl_set_of_books").Insert(&Html.Sob)
			if err != nil {
				if core.Conf.Debug {
					log.Println(err)
				}
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = "Fail to save the set of book"
			}
		}
	}
	if Html.Info.Show {
		err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_sobs_editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	} else {
		this.sobPage(w, r)
	}
}

func (this GLApp) sobEditPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title    string
		Sob      models.GLSetOfBook
		Timezone string
		Act      string
		Info     struct {
			Show    bool
			Type    string
			Message string
		}
	}
	Html.Timezone = core.Conf.TimeZone
	//Load account info;
	sobId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of gl.html?act=sobEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	ok, err := core.DBEngine.Table("gl_set_of_books").Where("id = ?", sobId).Get(&Html.Sob)
	if !ok {
		if err != nil {
			log.Println(err)
		}
		log.Printf("Illegal visit of gl.html?act=sobEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	Html.Title = "Edit set of book"
	Html.Act = "sobEdit"
	err = webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_sobs_editor.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) sobEdit(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title    string
		Sob      models.GLSetOfBook
		Timezone string
		Act      string
		Info     struct {
			Show    bool
			Type    string
			Message string
		}
	}
	Html.Timezone = core.Conf.TimeZone
	//Load account info;
	sobId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Illegal visit of gl.html?act=sobEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	ok, err := core.DBEngine.Table("gl_set_of_books").Where("id = ?", sobId).Get(&Html.Sob)
	if !ok {
		if err != nil {
			log.Println(err)
		}
		log.Printf("Illegal visit of gl.html?act=sobEdit")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	r.ParseForm()
	Html.Sob.SobName = r.Form.Get("sobname")
	Html.Sob.CurrencyCode = r.Form.Get("currencycode")
	Html.Sob.EolTime = core.UnixTime(r.Form.Get("eoltime"))
	Html.Title = "Edit set of book"
	Html.Act = "sobEdit"
	if Html.Sob.SobName == "" || Html.Sob.CurrencyCode == "" || Html.Sob.EolTime == 0 {
		Html.Info.Show = true
		Html.Info.Type = "danger"
		Html.Info.Message = "All input fields are required"
	}
	Html.Sob.LastUpdated = time.Now()
	if !Html.Info.Show {
		//Check conflict account;
		var dupSob models.GLSetOfBook
		ok, err := core.DBEngine.Table("gl_set_of_books").Where("sob_name=?", Html.Sob.SobName).And("id!=?", Html.Sob.Id).Get(&dupSob)
		if ok {
			Html.Info.Show = true
			Html.Info.Type = "danger"
			Html.Info.Message = "Found conflict account: "
			Html.Info.Message += dupSob.SobName
			Html.Info.Message += " ,currency code: "
			Html.Info.Message += dupSob.CurrencyCode
		} else {
			log.Println(Html.Sob)
			_, err = core.DBEngine.Table("gl_set_of_books").Where("id=?", sobId).Update(&Html.Sob)
			if err != nil {
				if core.Conf.Debug {
					log.Println(err)
				}
				Html.Info.Show = true
				Html.Info.Type = "danger"
				Html.Info.Message = "Fail to save the set of book"
			}
		}
	}
	if Html.Info.Show {
		err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_sobs_editor.html").Execute(w, Html)
		if err != nil {
			log.Println(err)
		}
	} else {
		this.sobPage(w, r)
	}
}

func (this GLApp) journalPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Sobs []models.GLSetOfBook
	}
	core.DBEngine.Table("gl_set_of_books").Where("1=1").Find(&Html.Sobs)
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_journals.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) journalAddPage(w http.ResponseWriter, r *http.Request) {
	var Html struct {
		Title    string
		Sobs     []models.GLSetOfBook
		Periods  []models.GLPeriod
		Accounts []models.GLAccount
		Journal  models.GLJournal
		Act      string
		Info     struct {
			Show    bool
			Type    string
			Message string
		}
	}
	core.DBEngine.Table("gl_set_of_books").Where("1=1").Find(&Html.Sobs)
	core.DBEngine.Table("gl_periods").Where("1=1").Find(&Html.Periods)
	core.DBEngine.Table("gl_accounts").Where("1=1").Find(&Html.Accounts)
	Html.Title = "Create Manual Journal"
	Html.Act = "journalAdd"
	err := webcore.GetTemplate(w, webcore.GetUILang(w, r), "gl_journals_editor.html").Execute(w, Html)
	if err != nil {
		log.Println(err)
	}
}

func (this GLApp) journalAdd(w http.ResponseWriter, r *http.Request) {
	strNullFail := `{"status":"failed","message":"Null input"}`
	strIllegalFail := `{"status":"failed","message":"Illegal input"}`
	strHeaderFail := `{"status":"failed","message":"Fail to save journal header"}`
	strLineFail := `{"status":"failed","message":"Fail to save journal line"}`
	strSuccess := `{"status":"success","message":"Operation completed"}`
	session := webcore.AclGetSession(r)
	if session == nil {
		w.Write([]byte(strIllegalFail))
		return
	}
	r.ParseForm()
	dataJson := r.Form.Get("data")
	if dataJson == "" {
		w.Write([]byte(strNullFail))
		return
	}
	var data struct {
		Header struct {
			SobId       string `json:"sob_id"`
			PeriodId    string `json:"period_id"`
			JournalDate string `json:"journal_date"`
			Voucher     string `json:"voucher"`
			Description string `json:"description"`
		} `json:"header"`
		Lines [][]string `json:"lines"`
	}
	err := json.Unmarshal([]byte(dataJson), &data)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
		}
		w.Write([]byte(strIllegalFail))
		return
	}
	//Insert header;
	glModel := new(models.GLModel)
	var header models.GLJournal
	header.Voucher = data.Header.Voucher
	if header.Voucher == "" {
		header.Voucher = glModel.NewVoucherNo()
	}
	header.Status = models.GLJournalStatusCreated
	header.Source = models.GLJournalSourceManual
	header.CreatedDate = core.NowTime()
	header.CreatedBy = session.User.NickName
	header.JournalDate = data.Header.JournalDate
	header.Description = data.Header.Description
	var nId int
	nId, err = strconv.Atoi(data.Header.SobId)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
		}
		w.Write([]byte(strIllegalFail))
		return
	}
	header.SobId = int64(nId)
	nId, err = strconv.Atoi(data.Header.PeriodId)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
		}
		w.Write([]byte(strIllegalFail))
		return
	}
	header.PeriodId = int64(nId)
	header.LastUpdated = time.Now()
	_, err = core.DBEngine.Table("gl_journals").Insert(&header)
	if err != nil {
		if core.Conf.Debug {
			log.Println(err)
		}
		w.Write([]byte(strHeaderFail))
		return
	}
	//Insert lines;
	for _, lData := range data.Lines {
		var line models.GLJournalEntry
		line.JournalId = header.Id
		line.Description = lData[0]
		nId, err = strconv.Atoi(lData[1])
		if err != nil {
			if core.Conf.Debug {
				log.Println(err)
			}
			w.Write([]byte(strIllegalFail))
			return
		}
		line.AccountId = int64(nId)
		var nMoney float64
		if lData[2] != "" {
			nMoney, err = strconv.ParseFloat(lData[2], 64)
			if err != nil {
				if core.Conf.Debug {
					log.Println(err)
				}
				nMoney = 0
			}
			line.Debit = float32(nMoney)
		}
		if lData[3] != "" {
			nMoney, err = strconv.ParseFloat(lData[3], 64)
			if err != nil {
				if core.Conf.Debug {
					log.Println(err)
				}
				nMoney = 0
			}
			line.Credit = float32(nMoney)
		}
		line.LastUpdated = time.Now()
		_, err := core.DBEngine.Table("gl_journal_entries").Insert(&line)
		if err != nil {
			if core.Conf.Debug {
				log.Println(line)
				log.Println(err)
			}
			w.Write([]byte(strLineFail))
			return
		}
	}
	w.Write([]byte(strSuccess))
}
