package email

type MailType string

const (
	MailTypeHtml  = MailType("html")
	MailTypePlain = MailType("plain")
)

type ContentConfig struct {
	Subject  string   `json:"subject" yaml:"subject"`
	Body     string   `json:"body" yaml:"body"`
	MailType MailType `json:"mail_type" yaml:"mail_type"`
	To       []string `json:"to" yaml:"to"`
	Cc       []string `json:"cc" yaml:"cc"`
}
