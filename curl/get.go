package curl

import (
	"encoding/json"
	"github.com/juxuny/data-utils/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func Get(url string, headers ...http.Header) (data []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}
	for _, h := range headers {
		for k, v := range h {
			if len(v) > 0 {
				for _, vo := range v {
					req.Header.Add(k, vo)
				}
			}
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf(resp.Status)
	}
	data, err = ioutil.ReadAll(resp.Body)
	return
}

func GetJson(url string, out interface{}, headers ...http.Header) error {
	data, err := Get(url, headers...)
	if err != nil {
		return err
	}
	log.Debug(string(data))
	return json.Unmarshal(data, out)
}
