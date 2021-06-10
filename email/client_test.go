package email

import (
	"testing"
)

func TestClient(t *testing.T) {
	c := NewClient(ClientConfig{
		User:        "fat-tiger@yandex.com",
		DisplayName: "fat-tiger",
		Password:    "",
		Host:        "smtp.yandex.com:465",
		Ssl:         true,
	})
	if err := c.Send(ContentConfig{
		Subject: "精选",
		Body: `真正让人变好的选择，过程都不会很舒服。你明知道躺在床上睡懒觉更舒服，但还是一早就起床；你明知道什么都不做比较轻松，但依旧选择追逐梦想。这就是生活，你必须坚持下去。
`,
		MailType: "plain",
		To: []string{
			"juxuny@163.com",
			"fat-tiger@neagk.com",
			"976813280@qq.com",
			"664841383@qq.com",
			"1031530366@qq.com",
		},
	}); err != nil {
		t.Fatal(err)
	}
}
