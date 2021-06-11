package proxy

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
