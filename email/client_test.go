package email

import (
	"github.com/gamexg/proxyclient"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	c := NewClient(ClientConfig{
		User:        "fat-tiger@yandex.com",
		DisplayName: "fat-tiger",
		Password:    "",
		Host:        "smtp.yandex.com:587",
		Ssl:         false,
	})
	if err := c.Send(ContentConfig{
		Subject: "真正让人变好的选择",
		Body: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>胖虎邀请你一起畅享生活</title>
    <style>
        body {
            margin: 0;
        }
    </style>
</head>
<body>
<img width="100%" height="auto" src="https://wx1.sinaimg.cn/mw2000/008ix0V3gy1graxvwxfw6j30dp0x1jxr.jpg">
</body>
</html>
`,
		MailType: "html",
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

func TestProxy(t *testing.T) {
	_ = os.Setenv("https_proxy", "http://127.0.0.1:7890")
	_ = os.Setenv("http_proxy", "http://127.0.0.1:7890")
	_ = os.Setenv("all_proxy", "socks5://127.0.0.1:7890")
	resp, err := http.Get("https://www.google.com")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	p, err := proxyclient.NewProxyClient("socks5://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}

	c, err := p.Dial("tcp", "www.google.com:80")
	if err != nil {
		panic(err)
	}

	io.WriteString(c, "GET / HTTP/1.0\r\nHOST:www.google.com\r\n\r\n")
	b, err := ioutil.ReadAll(c)
	if err != nil {
		panic(err)
	}
	t.Log(string(b))
}
