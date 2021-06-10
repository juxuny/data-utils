package email

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ProxyInfo struct {
	Ip   string
	Port int
}

func FetchProxyInfo(url string) (ret *ProxyInfo, err error) {
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
	ret = &ProxyInfo{
		Ip:   l[0],
		Port: int(port),
	}
	return ret, nil
}
