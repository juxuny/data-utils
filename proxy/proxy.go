package proxy

import (
	"fmt"
	"github.com/gamexg/proxyclient"
	"github.com/juxuny/data-utils/global_key"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/proxy/dt"
	"github.com/juxuny/env"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Ip   string
	Port int
}

type Auth proxy.Auth

func (t *Config) GetAddress(auth ...Auth) string {
	if len(auth) > 0 {
		return fmt.Sprintf("socks5://%s:%s@%s:%d", auth[0].User, auth[0].Password, t.Ip, t.Port)
	}
	return fmt.Sprintf("socks5://%s:%d", t.Ip, t.Port)
}

func FetchProxyInfo(url string) (ret *Config, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body failed")
	}
	address := string(data)
	l := strings.Split(address, ":")
	if len(l) != 2 {
		return nil, errors.Errorf("invalid address: %s", address)
	}
	port, err := strconv.ParseInt(strings.TrimSpace(l[1]), 10, 64)
	if err != nil {
		return nil, errors.Errorf("invalid port number: %v", l[1])
	}
	ret = &Config{
		Ip:   l[0],
		Port: int(port),
	}
	return ret, nil
}

type TestResult struct {
	ConnectDuration  time.Duration `json:"connect_duration"`
	SendDuration     time.Duration `json:"send_duration"`
	ResponseDuration time.Duration `json:"response_duration"`
	TotalLatency     time.Duration `json:"total_latency"`
	SuccessNum       int
	FailedNum        int
	ErrSummary       map[error]int
}

const DefaultTestTimeout = time.Second * 10

// address: proxy address, e.g socks5://user:pass@127.0.0.1:1080
func Test(address string, num int) (ret *TestResult, err error) {
	ret = &TestResult{}
	client, err := proxyclient.NewProxyClient(address)
	if err != nil {
		return nil, errors.Wrap(err, "create proxy client failed")
	}
	testHost := env.GetString(global_key.EnvKey.TestHost, "www.baidu.com:80")

	type innerTestResult struct {
		connectDuration, sendDuration, responseDuration, totalLatency time.Duration
	}

	var testFunc = func(client proxyclient.ProxyClient) (ret innerTestResult, err error) {
		start := time.Now()
		c, err := client.Dial("tcp", testHost)
		if err != nil {
			err = dt.ErrConnectFailed
			return
		}
		connectedTime := time.Now()
		ret.connectDuration = connectedTime.Sub(start)

		_, err = io.WriteString(c, fmt.Sprintf("GET / HTTP/1.0\r\nHOST:%s\r\n\r\n", testHost))
		if err != nil {
			err = dt.ErrSendDataFailed
			return
		}
		sentTime := time.Now()
		ret.sendDuration = sentTime.Sub(connectedTime)

		b, err := ioutil.ReadAll(c)
		if err != nil {
			err = dt.ErrReadDataFailed
			return
		}
		finishedTime := time.Now()
		ret.responseDuration = finishedTime.Sub(sentTime)
		if len(b) == 0 {
			err = dt.ErrEmptyResponse
			return
		}
		ret.totalLatency = finishedTime.Sub(start)
		return
	}
	for i := 0; i < num; i++ {
		middleResult, err := testFunc(client)
		if err != nil {
			ret.ErrSummary[err] += 1
			ret.FailedNum += 1
			continue
		}
		ret.ConnectDuration += middleResult.connectDuration
		ret.SendDuration += middleResult.sendDuration
		ret.ResponseDuration += middleResult.responseDuration
		ret.TotalLatency += middleResult.totalLatency
		ret.SuccessNum += 1
	}
	if ret.SuccessNum > 0 {
		ret.ConnectDuration /= time.Duration(ret.SuccessNum)
		ret.SendDuration /= time.Duration(ret.SuccessNum)
		ret.ResponseDuration /= time.Duration(ret.SuccessNum)
		ret.TotalLatency /= time.Duration(ret.SuccessNum)
	}

	return ret, nil
}

func RandServer(db *model.DB) (ret *model.Proxy, err error) {
	var count int
	scope := db.Where("latency > 0")
	if err := scope.Model(&model.Proxy{}).Count(&count).Error; err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "count failed")
	}
	if count == 0 {
		return nil, dt.ErrNotFound
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	offset := r.Intn(count)
	var out model.Proxy
	if err := scope.Offset(offset).First(&out).Error; err != nil {
		if !model.IsErrNoDataInDb(err) {
			log.Error(err)
			return nil, err
		}
	}
	return &out, nil
}
