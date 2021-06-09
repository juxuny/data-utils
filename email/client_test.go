package email

import (
	"github.com/juxuny/env"
	"testing"
)

func TestClient(t *testing.T) {
	c := NewClient(ClientConfig{
		User:        env.GetString("MAIL", "juxuny@163.com"),
		DisplayName: "fat-tiger",
		Password:    env.GetString("PASSWORD"),
		Host:        "smtp.163.com:465",
		Ssl:         true,
	})
	if err := c.Send(ContentConfig{
		Subject: "胖虎的美好生活",
		Body: `<!DOCTYPE html>
<html>
<head>
	<meta name=”viewport” content=”width=device-width, initial-scale=1, maximum-scale=1″>
	<meta charset="utf-8">
	<title></title>
	<style type="text/css">
		.box_style {
			margin: 20px;
			padding: 20px 10px;
			border-bottom: 1px dotted orange;
		}
		.text_style {
			font-weight: bold;
			font-size: 15px;
			width: 217px;
			text-align: left;
			display: inline-block;
		}
		.text_red_style {
			color: red;
		}
		.text_black_style {
			color: #000;
		}
	</style>
</head>
<body style="text-align: center;">
	<div class="box_style">
		<h2 style="text-align: center;">现在加入可享受如下特惠：</h2>
		<div style="text-align: center;">
			<span class="text_style text_black_style">1、<span class="text_red_style">85</span> zhe（电影票）</span>
		<br/>
		<span class="text_style text_black_style">2、每日大额外卖<span class="text_red_style">红包</span></span>
		<br/>
		<span class="text_style text_black_style">3、查询并领取各大电商<span class="text_red_style">券</span></span>
		<br/>
		<span class="text_style text_black_style">4、每次购物均可获取高额<span class="text_red_style">返利</span></span>
		<br/>
		<span class="text_style text_black_style">5、分享亦可<span class="text_red_style">赚</span></span>
		</div>
	</div>
	<div style="text-align: center;">
		<h3 style="text-align: center;">码在这,不显示的可以自行去<span class="text_red_style">vx</span>找:胖虎的变胖生活</h3>
		<img width="100px" height="100px" src="https://cdn.jsdelivr.net/gh/juxuny/res/qrcode/gh_eb717899ae68_430.jpeg">

	</div>
</body>
</html>
`,
		MailType: "html",
		To: []string{
			"juxuny@163.com",
			"fat-tiger@neagk.com",
			"976813280@qq.com",
			//"664841383@qq.com",
			//"1031530366@qq.com",
		},
	}); err != nil {
		t.Fatal(err)
	}
}
