package email

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
)

type ClientConfig struct {
	User        string `json:"user" yaml:"user"`
	DisplayName string `json:"display_name" yaml:"display_name"`
	Password    string `json:"password" yaml:"password"`
	Host        string `json:"host" yaml:"host"`
	Ssl         bool   `json:"ssl" yaml:"ssl"`
}

type client struct {
	config ClientConfig
}

func (t *client) send(user, sendUserName, password, host string, to []string, subject string, body string, mailType MailType) error {
	hp := strings.Split(host, ":")

	auth := smtp.PlainAuth("", user, password, hp[0])
	//var contentType string
	e := email.NewEmail()
	if mailType == MailTypeHtml {
		e.Headers.Add("ContentConfig-Type", "text/"+string(mailType)+"; charset=UTF-8")
		e.HTML = []byte(body)
	} else {
		e.Headers.Add("ContentConfig-Type", "text/plain; charset=UTF-8")
		e.Text = []byte(body)
	}
	e.To = to
	e.From = user
	e.Subject = subject
	if t.config.Ssl {
		return e.SendWithTLS(t.config.Host, auth, &tls.Config{ServerName: t.config.Host})
	} else {
		return e.Send(t.config.Host, auth)
	}
}

func (t *client) Send(content ContentConfig) error {
	return t.send(
		t.config.User,
		t.config.DisplayName,
		t.config.Password,
		t.config.Host,
		content.To,
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
