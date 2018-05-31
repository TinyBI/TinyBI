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

//Email queue for mailer;

import (
	"container/list"
	"crypto/tls"
	"log"
	"strconv"
	"time"
	"tinybi/core"
	"tinybi/logger"
	"tinybi/model"

	"github.com/go-gomail/gomail"
)

type EMail struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Contents    string
	Attachments []string
}

type mailEntry struct {
	Status    string
	LastError string
	Retries   uint
	Mail      EMail
}

const retriesMax uint = 3
const taskInterval time.Duration = 10

var queue *list.List

var running bool

var smtpInfo struct {
	Server string
	Port   int
	User   string
	Pass   string
	TLS    string
}

func init() {
	queue = list.New()
	smtpInfo.Server = "localhost"
	smtpInfo.Port = 25
}

func PushToQueue(mail EMail) {
	entry := mailEntry{Status: "PENDING", Retries: 0, Mail: mail}
	queue.PushBack(&entry)
}

func sendMailEry(mail *mailEntry) {
	goDail := gomail.NewDialer(smtpInfo.Server, smtpInfo.Port, smtpInfo.User, smtpInfo.Pass)
	if smtpInfo.TLS != "YES" {
		goDail.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	goMail := gomail.NewMessage()
	goMail.SetHeader("From", mail.Mail.From)
	goMail.SetHeader("To", mail.Mail.To...)
	goMail.SetHeader("Cc", mail.Mail.To...)
	goMail.SetHeader("Bcc", mail.Mail.To...)
	goMail.SetHeader("Subject", mail.Mail.Subject)
	goMail.SetBody("text/html", mail.Mail.Contents)
	//Attachments;
	for _, aPath := range mail.Mail.Attachments {
		goMail.Attach(aPath)
	}
	err := goDail.DialAndSend(goMail)
	if err != nil {
		mail.Status = "FAILED"
		mail.LastError = err.Error()
		logger.Printf("MAILER", "Fail to Send Mail, error:%s", mail.LastError)
		if core.Conf.Debug {
			log.Println("Fail to send mail", err)
		}
		return
	}
	mail.Status = "SUCCESS"
}

func runQueue() {
	for e := queue.Front(); e != nil; e = e.Next() {
		mail, ok := e.Value.(*mailEntry)
		if ok {
			switch mail.Status {
			case "PENDING":
				//Send the mail;
				sendMailEry(mail)
				break
			case "FAILED":
				//Retry it when retries is not reach the max value;
				if mail.Retries < retriesMax {
					mail.Retries++
					sendMailEry(mail)
				} else {
					queue.Remove(e)
				}
				break
			case "SUCCESS":
				//Already sent, remove it from queue;
				queue.Remove(e)
				break
			}
		}
	}
}

func run() {
	for running {
		select {
		case <-time.After(taskInterval * time.Second):
			runQueue()
			break
		}
	}
}

func Start() {
	//Read SMTP configuration from settings;
	smtpServer := model.BusinessSettings.Get("SMTP_SERVER")
	if smtpServer != nil {
		smtpInfo.Server = smtpServer.Value
	}
	smtpPort := model.BusinessSettings.Get("SMTP_PORT")
	if smtpPort != nil {
		var err error
		smtpInfo.Port, err = strconv.Atoi(smtpPort.Value)
		if err != nil {
			smtpInfo.Port = 25
		}
	}
	smtpUser := model.BusinessSettings.Get("SMTP_USER")
	if smtpServer != nil {
		smtpInfo.User = smtpUser.Value
	}
	smtpPass := model.BusinessSettings.Get("SMTP_PASS")
	if smtpServer != nil {
		smtpInfo.Pass = smtpPass.Value
	}
	smtpTLS := model.BusinessSettings.Get("SMTP_TLS")
	if smtpTLS != nil {
		smtpInfo.TLS = smtpTLS.Value
	}
	running = true
	go run()
}

func Stop() {
	running = false
}
