package email

import "testing"

func TestClient(t *testing.T) {
	c := NewClient(ClientConfig{
		User:        "ykeoh9357@163.com",
		DisplayName: "fat-tiger",
		Password:    "",
		Host:        "smtp.163.com:25",
		Ssl:         false,
	})
	if err := c.Send(ContentConfig{
		Subject: "胖虎的美好生成",
		Body: `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="iso-8859-15">
			<title>MMOGA POWER</title>
		</head>
		<body>
			GO 发送邮件，官方连包都帮我们写好了，真是贴心啊！！！
		</body>
		</html>`,
		MailType: "",
		To: []string{
			"juxuny@gmail.com",
			"fat-tiger@neagk.com",
		},
	}); err != nil {
		t.Fatal(err)
	}
}
