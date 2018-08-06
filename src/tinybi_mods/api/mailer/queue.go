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
package mailer

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	coreMailer "tinybi/mailer"
	"strings"
)

type MailInfo struct {
	Recipients  string `json "recipients"`
	Cc          string `json "cc"`
	Bcc         string `json "bcc"`
	Subject     string `json "subject"`
	Contents    string `json "contents"`
	Attachments string `json "attachments"`
}

func Post(w http.ResponseWriter, r *http.Request) (bool, string, interface{}) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, "Empty request", nil
	}
	var mailInfo MailInfo
	err = json.Unmarshal(body, &mailInfo)
	if err != nil {
		return false, "Wrong data format", body
	}
	if mailInfo.Subject == "" {
		return false, "Empty Subject", body
	}
	if mailInfo.Recipients == "" {
		return false, "Empty Recipient", body
	}
	var mail coreMailer.EMail
	mail.To = strings.Split(mailInfo.Recipients, ",")
	if mailInfo.Cc != "" {
		mail.Cc = strings.Split(mailInfo.Cc, ",")
	}
	if mailInfo.Bcc != "" {
		mail.Bcc = strings.Split(mailInfo.Bcc, ",")
	}
	mail.Subject = mailInfo.Subject
	mail.Contents = mailInfo.Contents
	if mailInfo.Attachments != "" {
		mail.Attachments = strings.Split(mailInfo.Attachments, ",")
	}
	coreMailer.PushToQueue(mail)
	return true, "", nil
}
