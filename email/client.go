package email

import (
	"fmt"
	"github.com/juxuny/gomail"
	"strconv"
	"strings"
	"time"
)

type ClientConfig struct {
	User        string   `json:"user" yaml:"user"`
	DisplayName string   `json:"display_name" yaml:"display_name"`
	Password    string   `json:"password" yaml:"password"`
	Host        string   `json:"host" yaml:"host"`
	Ssl         bool     `json:"ssl" yaml:"ssl"`
	CC          []string `json:"cc" yaml:"cc"`
}

type client struct {
	config ClientConfig
}

//func (t *client) send(user, sendUserName, password, host string, to []string, subject string, body string, mailType MailType) error {
//	hp := strings.Split(host, ":")
//
//	auth := smtp.PlainAuth("", user, password, hp[0])
//	//var contentType string
//	e := email.NewEmail()
//	if mailType == MailTypeHtml {
//		e.Headers.Add("ContentConfig-Type", "text/"+string(mailType)+"; charset=UTF-8")
//		e.HTML = []byte(body)
//	} else {
//		e.Headers.Add("ContentConfig-Type", "text/plain; charset=UTF-8")
//		e.Text = []byte(body)
//	}
//	e.To = to
//	e.From = user
//	e.Subject = subject
//	if t.config.Ssl {
//		return e.SendWithTLS(t.config.Host, auth, &tls.Config{ServerName: hp[0]})
//	} else {
//		return e.Send(t.config.Host, auth)
//	}
//}

func (t *client) send(user, sendUserName, password, host string, to []string, cc []string, subject string, body string, mailType MailType) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(user, sendUserName))
	//m.SetHeader("To", to...)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetDateHeader("Date", time.Now())
	toMap := map[string]bool{}
	for _, item := range to {
		toMap[item] = true
	}
	finalCc := make([]string, 0)
	if len(t.config.CC) > 0 {
		for _, item := range t.config.CC {
			if !toMap[item] {
				finalCc = append(finalCc, item)
			}
		}
	}
	if len(cc) > 0 {
		for _, item := range cc {
			if !toMap[item] {
				finalCc = append(finalCc, item)
			}
		}
	}
	m.SetHeader("Cc", finalCc...)
	m.SetHeader("Message-ID", fmt.Sprintf("%v", time.Now().UnixNano()))
	//m.AddAlternative("text/plain", "hello", gomail.SetPartEncoding(gomail.Base64))
	m.SetBody("text/"+string(mailType), body, gomail.SetPartEncoding(gomail.Base64))
	hp := strings.Split(host, ":")
	var port int64
	if len(hp) == 1 {
		port = 25
	} else {
		tmp, err := strconv.ParseInt(hp[1], 10, 64)
		if err != nil {
			return err
		}
		port = tmp
	}
	d := gomail.NewDialer(hp[0], int(port), user, password)
	d.SSL = t.config.Ssl
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (t *client) Send(content ContentConfig) error {
	return t.send(
		t.config.User,
		t.config.DisplayName,
		t.config.Password,
		t.config.Host,
		content.To,
		content.Cc,
		content.Subject,
		content.Body,
		content.MailType,
	)
}

func NewClient(config ClientConfig) *client {
	c := &client{
		config: config,
	}
	return c
}
