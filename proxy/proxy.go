package proxy

import (
	"context"
	"fmt"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/proxy/dt"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
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
	ret = &TestResult{
		ErrSummary: make(map[error]int),
	}
	//client, err := proxyclient.NewProxyClient(address)
	//if err != nil {
	//	return nil, errors.Wrap(err, "create proxy client failed")
	//}
	//testHost := env.GetString(global_key.EnvKey.TestHost, "www.baidu.com:80")

	type innerTestResult struct {
		connectDuration, sendDuration, responseDuration, totalLatency time.Duration
	}

	var testFunc = func() (ret innerTestResult, err error) {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), DefaultTestTimeout)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.baidu.com", nil)
		if err != nil {
			return ret, err
		}
		proxyUrl, err := url.Parse(address)
		httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		httpClient.Timeout = DefaultTestTimeout
		resp, err := httpClient.Do(req)
		if err != nil {
			return ret, dt.ErrSendDataFailed
		}
		if resp.StatusCode/100 != 2 {
			return ret, errors.Errorf(resp.Status)
		}
		sentTime := time.Now()
		ret.sendDuration = start.Sub(sentTime)
		b, err := ioutil.ReadAll(resp.Body)
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
		middleResult, err := testFunc()
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
