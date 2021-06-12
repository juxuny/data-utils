package proxy

import (
	"testing"
)

func TestTest(t *testing.T) {
	testAddress := "socks5://127.0.0.1:7890"
	ret, err := Test(testAddress, 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("connect_duration: ", ret.ConnectDuration)
	t.Log("send_duration: ", ret.SendDuration)
	t.Log("response_duration: ", ret.ResponseDuration)
	t.Log("total_latency: ", ret.TotalLatency)
}
